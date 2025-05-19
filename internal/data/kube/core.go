package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type core struct {
	kubeMap KubeMap
}

func NewCore(kubeMap KubeMap) service.KubeCore {
	return &core{
		kubeMap: kubeMap,
	}
}

var _ service.KubeCore = (*core)(nil)

func (r *core) GetNamespace(ctx context.Context, uuid, facility, name string) (*v1.Namespace, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.CoreV1().Namespaces().Get(ctx, name, opts)
}

func (r *core) CreateNamespace(ctx context.Context, uuid, facility, name string) (*v1.Namespace, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	opts := metav1.CreateOptions{}
	return clientset.CoreV1().Namespaces().Create(ctx, namespace, opts)
}

func (r *core) ListServices(ctx context.Context, uuid, facility, namespace string) ([]v1.Service, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
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

func (r *core) ListPods(ctx context.Context, uuid, facility, namespace string) ([]v1.Pod, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
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

func (r *core) ListPersistentVolumeClaims(ctx context.Context, uuid, facility, namespace string) ([]v1.PersistentVolumeClaim, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
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
