package kube

import (
	"bytes"
	"context"
	"io"
	"net/url"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"

	oscore "github.com/otterscale/otterscale/internal/core"
)

type core struct {
	kube *Kube
}

func NewCore(kube *Kube) oscore.KubeCoreRepo {
	return &core{
		kube: kube,
	}
}

var _ oscore.KubeCoreRepo = (*core)(nil)

func (r *core) GetNode(ctx context.Context, config *rest.Config, name string) (*oscore.Node, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.CoreV1().Nodes().Get(ctx, name, opts)
}

func (r *core) UpdateNode(ctx context.Context, config *rest.Config, node *oscore.Node) (*oscore.Node, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.UpdateOptions{}
	return clientset.CoreV1().Nodes().Update(ctx, node, opts)
}

func (r *core) ListNamespaces(ctx context.Context, config *rest.Config) ([]corev1.Namespace, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.CoreV1().Namespaces().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *core) ListServices(ctx context.Context, config *rest.Config, namespace string) ([]oscore.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.CoreV1().Services(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) ListServicesByOptions(ctx context.Context, config *rest.Config, namespace, label, field string) ([]oscore.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}
	list, err := clientset.CoreV1().Services(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) ListVirtualMachineServices(ctx context.Context, config *rest.Config, namespace, vmName string) ([]oscore.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	if vmName != "" {
		opts.LabelSelector = oscore.VirtualMachineNameLabel + "=" + vmName
	}
	list, err := clientset.CoreV1().Services(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) GetService(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func (r *core) CreateVirtualMachineService(ctx context.Context, config *rest.Config, namespace, name, vmName string, ports []corev1.ServicePort) (*oscore.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				oscore.VirtualMachineNameLabel: vmName,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
			Type:  corev1.ServiceTypeNodePort,
		},
	}
	opts := metav1.CreateOptions{}
	return clientset.CoreV1().Services(namespace).Create(ctx, service, opts)
}

func (r *core) UpdateService(ctx context.Context, config *rest.Config, namespace string, service *corev1.Service) (*oscore.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.UpdateOptions{}
	return clientset.CoreV1().Services(namespace).Update(ctx, service, opts)
}

func (r *core) DeleteService(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.CoreV1().Services(namespace).Delete(ctx, name, opts)
}

func (r *core) ListPods(ctx context.Context, config *rest.Config, namespace string) ([]oscore.Pod, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) ListPodsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]oscore.Pod, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: label,
	}
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) DeletePod(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}

func (r *core) GetLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (string, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return "", err
	}

	opts := corev1.PodLogOptions{
		Container: containerName,
	}
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &opts)

	logStream, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer logStream.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logStream)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (r *core) StreamLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (io.ReadCloser, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := corev1.PodLogOptions{
		Container: containerName,
		Follow:    true,
	}
	return clientset.CoreV1().Pods(namespace).GetLogs(podName, &opts).Stream(ctx)
}

func (r *core) CreateExecutor(config *rest.Config, namespace, podName, containerName string, command []string) (remotecommand.Executor, error) {
	clientset, err := r.kube.clientset(config)
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

	return r.createExecutor(config, req.URL())
}

func (r *core) ListPersistentVolumeClaims(ctx context.Context, config *rest.Config, namespace string) ([]oscore.PersistentVolumeClaim, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) GetPersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.PersistentVolumeClaim, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, opts)
}

func (r *core) PatchPersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string, data []byte) (*oscore.PersistentVolumeClaim, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.PatchOptions{}
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Patch(ctx, name, types.JSONPatchType, data, opts)
}

func (r *core) GetNamespace(ctx context.Context, config *rest.Config, name string) (*oscore.Namespace, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.CoreV1().Namespaces().Get(ctx, name, opts)
}

func (r *core) CreateNamespace(ctx context.Context, config *rest.Config, name string) (*oscore.Namespace, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	opts := metav1.CreateOptions{}
	return clientset.CoreV1().Namespaces().Create(ctx, namespace, opts)
}

func (r *core) GetConfigMap(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.ConfigMap, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, opts)
}

func (r *core) CreateConfigMap(ctx context.Context, config *rest.Config, namespace, name string, data map[string]string) (*oscore.ConfigMap, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	opts := metav1.CreateOptions{}
	return clientset.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, opts)
}

func (r *core) GetSecret(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.Secret, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.CoreV1().Secrets(namespace).Get(ctx, name, opts)
}

// https://github.com/kubernetes/kubectl/blob/45c6a75b21af19de57b586862dc509a5d7afc081/pkg/cmd/exec/exec.go#L145
func (r *core) createExecutor(config *rest.Config, url *url.URL) (remotecommand.Executor, error) {
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", url)
	if err != nil {
		return nil, err
	}
	websocketExec, err := remotecommand.NewWebSocketExecutor(config, "GET", url.String())
	if err != nil {
		return nil, err
	}
	exec, err = remotecommand.NewFallbackExecutor(websocketExec, exec, func(err error) bool {
		return httpstream.IsUpgradeFailure(err) || httpstream.IsHTTPSProxyError(err)
	})
	if err != nil {
		return nil, err
	}
	return exec, nil
}
