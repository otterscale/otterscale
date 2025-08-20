package kube

import (
	"bytes"
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"

	oscore "github.com/openhdc/otterscale/internal/core"
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

func (r *core) CreateService(ctx context.Context, config *rest.Config, namespace, name string, spec *corev1.ServiceSpec) (*oscore.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	if spec != nil {
		svc.Spec = *spec
	}

	opts := metav1.CreateOptions{}
	return clientset.CoreV1().Services(namespace).Create(ctx, svc, opts)
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

func (r *core) UpdateService(ctx context.Context, config *rest.Config, namespace, name string, spec *corev1.ServiceSpec) (*oscore.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	// Get existing to keep resourceVersion, immutable fields, etc.
	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if spec != nil {
		svc.Spec = *spec
	}

	opts := metav1.UpdateOptions{}
	return clientset.CoreV1().Services(namespace).Update(ctx, svc, opts)
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

func (r *core) GetPodLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (string, error) {
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

func (r *core) CreatePersistentVolumeClaims(ctx context.Context, config *rest.Config, namespace, name string, spec *corev1.PersistentVolumeClaimSpec) (*oscore.PersistentVolumeClaim, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	if spec != nil {
		pvc.Spec = *spec
	}

	opts := metav1.CreateOptions{}
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, opts)
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

func (r *core) UpdatePersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string, spec *corev1.PersistentVolumeClaimSpec) (*oscore.PersistentVolumeClaim, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	if spec != nil {
		pvc.Spec = *spec
	}

	opts := metav1.UpdateOptions{}
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, pvc, opts)
}

func (r *core) DeletePersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, opts)
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
