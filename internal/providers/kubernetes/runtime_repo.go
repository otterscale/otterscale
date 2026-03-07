package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/transport/spdy"

	"github.com/otterscale/otterscale/internal/core"
)

// runtimeRepo implements core.RuntimeRepo by delegating to the
// Kubernetes typed, dynamic, and SPDY clients, accessed through
// the tunnel.
type runtimeRepo struct {
	kubernetes *Kubernetes
}

// NewRuntimeRepo returns a core.RuntimeRepo backed by Kubernetes.
func NewRuntimeRepo(kubernetes *Kubernetes) core.RuntimeRepo {
	return &runtimeRepo{kubernetes: kubernetes}
}

var _ core.RuntimeRepo = (*runtimeRepo)(nil)

// ---------------------------------------------------------------------------
// PodLogs
// ---------------------------------------------------------------------------

// PodLogs opens a streaming log reader for a container.
func (r *runtimeRepo) PodLogs(ctx context.Context, cluster, namespace, name string, opts core.PodLogOptions) (io.ReadCloser, error) {
	clientset, err := r.clientset(ctx, cluster)
	if err != nil {
		return nil, err
	}

	logOpts := &corev1.PodLogOptions{
		Container:  opts.Container,
		Follow:     opts.Follow,
		Previous:   opts.Previous,
		Timestamps: opts.Timestamps,
	}
	if opts.TailLines != nil {
		logOpts.TailLines = opts.TailLines
	}
	if opts.SinceSeconds != nil {
		logOpts.SinceSeconds = opts.SinceSeconds
	}
	if opts.SinceTime != nil {
		logOpts.SinceTime = &metav1.Time{Time: *opts.SinceTime}
	}
	if opts.LimitBytes != nil {
		logOpts.LimitBytes = opts.LimitBytes
	}

	result, err := clientset.CoreV1().Pods(namespace).GetLogs(name, logOpts).Stream(ctx)
	return result, wrapK8sError(err)
}

// ---------------------------------------------------------------------------
// Exec
// ---------------------------------------------------------------------------

// Exec starts an interactive exec session and blocks until it completes.
func (r *runtimeRepo) Exec(ctx context.Context, cluster, namespace, name string, opts *core.ExecOptions) error {
	config, err := r.kubernetes.spdyConfig(ctx, cluster)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &core.DomainError{Code: core.ErrorCodeInternal, Message: "create clientset for exec", Cause: err}
	}

	execOpts := &corev1.PodExecOptions{
		Container: opts.Container,
		Command:   opts.Command,
		TTY:       opts.TTY,
		Stdin:     opts.Stdin != nil,
		Stdout:    opts.Stdout != nil,
		Stderr:    opts.Stderr != nil,
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(name).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(execOpts, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
	if err != nil {
		return &core.DomainError{Code: core.ErrorCodeInternal, Message: "create SPDY executor", Cause: err}
	}

	streamOpts := remotecommand.StreamOptions{
		Stdin:  opts.Stdin,
		Stdout: opts.Stdout,
		Stderr: opts.Stderr,
		Tty:    opts.TTY,
	}
	if opts.TTY && opts.SizeQueue != nil {
		streamOpts.TerminalSizeQueue = &sizeQueueAdapter{inner: opts.SizeQueue}
	}

	return wrapK8sError(executor.StreamWithContext(ctx, streamOpts))
}

// ---------------------------------------------------------------------------
// Scale
// ---------------------------------------------------------------------------

// GetScale reads the current replica count via the /scale subresource.
func (r *runtimeRepo) GetScale(ctx context.Context, cluster string, gvr schema.GroupVersionResource, namespace, name string) (int32, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return 0, err
	}

	scaleObj, err := client.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{}, "scale")
	if err != nil {
		return 0, wrapK8sError(err)
	}

	replicas, found, err := unstructured.NestedInt64(scaleObj.Object, "spec", "replicas")
	if err != nil || !found {
		return 0, &core.DomainError{Code: core.ErrorCodeInternal, Message: "failed to read spec.replicas from scale subresource"}
	}

	if replicas < math.MinInt32 || replicas > math.MaxInt32 {
		return 0, &core.DomainError{Code: core.ErrorCodeInternal, Message: "spec.replicas out of int32 range"}
	}
	return int32(replicas), nil
}

// UpdateScale sets the desired replica count via the /scale subresource.
func (r *runtimeRepo) UpdateScale(ctx context.Context, cluster string, gvr schema.GroupVersionResource, namespace, name string, replicas int32) (int32, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return 0, err
	}

	// GET current scale
	scaleObj, err := client.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{}, "scale")
	if err != nil {
		return 0, wrapK8sError(err)
	}

	// SET desired replicas
	if err := unstructured.SetNestedField(scaleObj.Object, int64(replicas), "spec", "replicas"); err != nil {
		return 0, &core.DomainError{Code: core.ErrorCodeInternal, Message: "set spec.replicas", Cause: err}
	}

	// UPDATE scale subresource
	updated, err := client.Resource(gvr).Namespace(namespace).Update(ctx, scaleObj, metav1.UpdateOptions{}, "scale")
	if err != nil {
		return 0, wrapK8sError(err)
	}

	newReplicas, found, err := unstructured.NestedInt64(updated.Object, "spec", "replicas")
	if err != nil {
		return 0, &core.DomainError{Code: core.ErrorCodeInternal, Message: "read updated replicas", Cause: err}
	}
	if !found {
		return 0, &core.DomainError{Code: core.ErrorCodeInternal, Message: "spec.replicas not found in updated scale subresource"}
	}
	if newReplicas < math.MinInt32 || newReplicas > math.MaxInt32 {
		return 0, &core.DomainError{Code: core.ErrorCodeInternal, Message: "updated spec.replicas out of int32 range"}
	}
	return int32(newReplicas), nil
}

// ---------------------------------------------------------------------------
// Restart
// ---------------------------------------------------------------------------

// Restart triggers a rolling restart by patching the pod template
// annotation with kubectl.kubernetes.io/restartedAt, equivalent to
// `kubectl rollout restart`.
func (r *runtimeRepo) Restart(ctx context.Context, cluster string, gvr schema.GroupVersionResource, namespace, name string) error {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return err
	}

	// time.Now is used directly (not injected) because the annotation
	// value only needs to differ from the previous value to trigger a
	// rolling update — its exact timestamp is not significant for
	// correctness or testability.
	patchData := map[string]any{
		"spec": map[string]any{
			"template": map[string]any{
				"metadata": map[string]any{
					"annotations": map[string]any{
						"kubectl.kubernetes.io/restartedAt": time.Now().UTC().Format(time.RFC3339),
					},
				},
			},
		},
	}
	data, err := json.Marshal(patchData)
	if err != nil {
		return fmt.Errorf("marshal restart patch: %w", err)
	}

	_, err = client.Resource(gvr).Namespace(namespace).Patch(ctx, name, types.MergePatchType, data, metav1.PatchOptions{})
	return wrapK8sError(err)
}

// ---------------------------------------------------------------------------
// PortForward
// ---------------------------------------------------------------------------

// PortForward opens a port-forward session via SPDY and copies data
// bidirectionally between the caller's stdin/stdout and the pod.
// It waits for both copy directions to complete before returning.
func (r *runtimeRepo) PortForward(ctx context.Context, cluster, namespace, name string, opts core.PortForwardOptions) error {
	config, err := r.kubernetes.spdyConfig(ctx, cluster)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &core.DomainError{Code: core.ErrorCodeInternal, Message: "create clientset for port-forward", Cause: err}
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(name).
		Namespace(namespace).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return &core.DomainError{Code: core.ErrorCodeInternal, Message: "create SPDY round-tripper", Cause: err}
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())
	streamConn, _, err := dialer.Dial(portForwardProtocolV1)
	if err != nil {
		return wrapK8sError(err)
	}
	defer streamConn.Close()

	return r.runPortForwardStreams(ctx, streamConn, opts)
}

// runPortForwardStreams creates SPDY streams for the port-forward
// session and copies data bidirectionally. It is extracted from
// PortForward to reduce function length.
func (r *runtimeRepo) runPortForwardStreams(ctx context.Context, streamConn httpstream.Connection, opts core.PortForwardOptions) error {
	portStr := strconv.FormatInt(int64(opts.Port), 10)
	requestID := "0"

	// Create error stream.
	errorHeaders := http.Header{}
	errorHeaders.Set(corev1.StreamType, corev1.StreamTypeError)
	errorHeaders.Set(corev1.PortHeader, portStr)
	errorHeaders.Set(corev1.PortForwardRequestIDHeader, requestID)

	errorStream, err := streamConn.CreateStream(errorHeaders)
	if err != nil {
		return &core.DomainError{Code: core.ErrorCodeInternal, Message: "create error stream", Cause: err}
	}
	defer errorStream.Close()

	// Create data stream.
	dataHeaders := http.Header{}
	dataHeaders.Set(corev1.StreamType, corev1.StreamTypeData)
	dataHeaders.Set(corev1.PortHeader, portStr)
	dataHeaders.Set(corev1.PortForwardRequestIDHeader, requestID)

	dataStream, err := streamConn.CreateStream(dataHeaders)
	if err != nil {
		return &core.DomainError{Code: core.ErrorCodeInternal, Message: "create data stream", Cause: err}
	}
	defer dataStream.Close()

	var wg sync.WaitGroup

	// Check for immediate errors from kubelet.
	wg.Add(1)
	go func() {
		defer wg.Done()
		const errorBufSize = 1024
		buf := make([]byte, errorBufSize)
		n, _ := errorStream.Read(buf)
		if n > 0 {
			if err := dataStream.Close(); err != nil {
				slog.Warn("failed to close data stream after kubelet error", "error", err)
			}
		}
	}()

	// Bidirectional copy — wait for BOTH directions to complete.
	errCh := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(dataStream, opts.Stdin)
		errCh <- err
	}()

	go func() {
		defer wg.Done()
		_, err := io.Copy(opts.Stdout, dataStream)
		errCh <- err
	}()

	var firstErr error
	for i := 0; i < 2; i++ {
		select {
		case <-ctx.Done():
			streamConn.Close()
			wg.Wait()
			return ctx.Err()
		case err := <-errCh:
			if err != nil && firstErr == nil {
				firstErr = err
				streamConn.Close()
			}
		}
	}

	wg.Wait()
	return firstErr
}

// portForwardProtocolV1 is the subprotocol used for Kubernetes port
// forwarding over SPDY.
const portForwardProtocolV1 = "portforward.k8s.io"

// ---------------------------------------------------------------------------
// SubResourceAction
// ---------------------------------------------------------------------------

// SubResourceAction invokes a PUT or POST action on a named
// subresource. This covers use-cases like KubeVirt VM
// start/stop/restart/migrate where the API server exposes
// state-transition endpoints as subresources.
func (r *runtimeRepo) SubResourceAction(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
	namespace, name, subresource, method string, body []byte,
) (map[string]any, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return nil, err
	}

	obj := &unstructured.Unstructured{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &obj.Object); err != nil {
			return nil, &core.DomainError{Code: core.ErrorCodeInvalidArgument, Message: "invalid JSON body", Cause: err}
		}
	}
	if obj.Object == nil {
		obj.Object = map[string]any{}
	}
	obj.SetName(name)
	obj.SetNamespace(namespace)

	var result *unstructured.Unstructured
	switch method {
	case "PUT":
		result, err = client.Resource(gvr).Namespace(namespace).Update(ctx, obj, metav1.UpdateOptions{}, subresource)
	case "POST":
		result, err = client.Resource(gvr).Namespace(namespace).Create(ctx, obj, metav1.CreateOptions{}, subresource)
	default:
		return nil, &core.DomainError{Code: core.ErrorCodeInvalidArgument, Message: fmt.Sprintf("unsupported method %q for subresource action", method)}
	}
	if err != nil {
		return nil, wrapK8sError(err)
	}
	if result == nil {
		return nil, nil
	}
	return result.Object, nil
}

// ---------------------------------------------------------------------------
// Terminal size adapter
// ---------------------------------------------------------------------------

// sizeQueueAdapter bridges the domain core.TerminalSizer interface to
// the remotecommand.TerminalSizeQueue interface required by SPDY
// executors. This keeps the domain layer free of client-go dependencies.
type sizeQueueAdapter struct {
	inner core.TerminalSizer
}

func (a *sizeQueueAdapter) Next() *remotecommand.TerminalSize {
	s := a.inner.Next()
	if s == nil {
		return nil
	}
	return &remotecommand.TerminalSize{Width: s.Width, Height: s.Height}
}

// ---------------------------------------------------------------------------
// Client helpers
// ---------------------------------------------------------------------------

// clientset builds a fresh impersonated typed Kubernetes clientset for
// the given cluster. A new clientset is created per request because
// each request may carry different impersonation credentials (user
// subject + groups). The underlying HTTP transport is cached
// per-cluster in Kubernetes.roundTripper, so only the Go-level wrapper
// is allocated — negligible compared to the actual API call latency.
func (r *runtimeRepo) clientset(ctx context.Context, cluster string) (*kubernetes.Clientset, error) {
	config, err := r.kubernetes.impersonationConfig(ctx, cluster)
	if err != nil {
		return nil, err
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, &core.DomainError{Code: core.ErrorCodeInternal, Message: "create kubernetes clientset", Cause: err}
	}
	return cs, nil
}

// dynamicClient builds a fresh impersonated dynamic client for the
// given cluster. See clientset for the rationale on per-request
// client creation.
func (r *runtimeRepo) dynamicClient(ctx context.Context, cluster string) (*dynamic.DynamicClient, error) {
	config, err := r.kubernetes.impersonationConfig(ctx, cluster)
	if err != nil {
		return nil, err
	}
	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, &core.DomainError{Code: core.ErrorCodeInternal, Message: "create dynamic client for runtime", Cause: err}
	}
	return dc, nil
}
