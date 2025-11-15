package kubernetes

import (
	"context"
	"io"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/otterscale/kubevirt-client-go/containerizeddataimporter/scheme"

	"github.com/otterscale/otterscale/internal/core/application/workload"
)

type podRepo struct {
	kubernetes *Kubernetes
}

func NewPodRepo(kubernetes *Kubernetes) workload.PodRepo {
	return &podRepo{
		kubernetes: kubernetes,
	}
}

var _ workload.PodRepo = (*podRepo)(nil)

func (r *podRepo) List(ctx context.Context, scope, namespace, selector string) ([]workload.Pod, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CoreV1().Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *podRepo) Get(ctx context.Context, scope, namespace, name string) (*workload.Pod, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CoreV1().Pods(namespace).Get(ctx, name, opts)
}

func (r *podRepo) Create(ctx context.Context, scope, namespace string, p *workload.Pod) (*workload.Pod, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.CoreV1().Pods(namespace).Create(ctx, p, opts)
}

func (r *podRepo) Update(ctx context.Context, scope, namespace string, p *workload.Pod) (*workload.Pod, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.CoreV1().Pods(namespace).Update(ctx, p, opts)
}

func (r *podRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}

func (r *podRepo) Stream(ctx context.Context, scope, namespace, podName, containerName string, duration time.Duration, follow bool) (io.ReadCloser, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := corev1.PodLogOptions{
		Container: containerName,
		Follow:    follow,
	}

	if duration > 0 {
		sec := int64(duration.Round(time.Second).Seconds())
		opts.SinceSeconds = &sec
	}

	return clientset.CoreV1().Pods(namespace).GetLogs(podName, &opts).Stream(ctx)
}

func (r *podRepo) Executer(scope, namespace, podName, containerName string, command []string) (remotecommand.Executor, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	config, err := r.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	// https://github.com/kubernetes/kubectl/blob/45c6a75b21af19de57b586862dc509a5d7afc081/pkg/cmd/exec/exec.go#L385
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
		Container: containerName,
		Command:   command,
	}, scheme.ParameterCodec)

	// https://github.com/kubernetes/kubectl/blob/45c6a75b21af19de57b586862dc509a5d7afc081/pkg/cmd/exec/exec.go#L145

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, err
	}

	websocketExec, err := remotecommand.NewWebSocketExecutor(config, "GET", req.URL().String())
	if err != nil {
		return nil, err
	}

	return remotecommand.NewFallbackExecutor(websocketExec, exec, func(err error) bool {
		return httpstream.IsUpgradeFailure(err) || httpstream.IsHTTPSProxyError(err)
	})
}
