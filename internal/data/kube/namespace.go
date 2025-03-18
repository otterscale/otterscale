package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type namespace struct {
	kubes Kubes
}

func NewNamespace(kubes Kubes) service.KubeNamespace {
	return &namespace{
		kubes: kubes,
	}
}

var _ service.KubeNamespace = (*namespace)(nil)

func (r *namespace) Get(ctx context.Context, cluster, name string) (*corev1.Namespace, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return client.CoreV1().Namespaces().Get(ctx, name, opts)
}

func (r *namespace) Create(ctx context.Context, cluster, name string) (*corev1.Namespace, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	opts := metav1.CreateOptions{}
	return client.CoreV1().Namespaces().Create(ctx, namespace, opts)
}
