package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type core struct {
	kubes Kubes
}

func NewCore(kubes Kubes) service.KubeCore {
	return &core{
		kubes: kubes,
	}
}

var _ service.KubeCore = (*core)(nil)

func (r *core) GetNamespace(ctx context.Context, cluster, name string) (*v1.Namespace, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return client.CoreV1().Namespaces().Get(ctx, name, opts)
}

func (r *core) CreateNamespace(ctx context.Context, cluster, name string) (*v1.Namespace, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	opts := metav1.CreateOptions{}
	return client.CoreV1().Namespaces().Create(ctx, namespace, opts)
}

func (r *core) ListServices(ctx context.Context, cluster, namespace string) ([]v1.Service, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := client.CoreV1().Services(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) ListPods(ctx context.Context, cluster, namespace string) ([]v1.Pod, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := client.CoreV1().Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) ListPersistentVolumeClaims(ctx context.Context, cluster, namespace string) ([]v1.PersistentVolumeClaim, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
