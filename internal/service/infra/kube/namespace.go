package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openhdc/openhdc/internal/service/domain/service"
)

type namespace struct {
	client *kubernetes.Clientset
}

func NewNamespace(client *kubernetes.Clientset) service.KubeNamespace {
	return &namespace{
		client: client,
	}
}

var _ service.KubeNamespace = (*namespace)(nil)

func (r *namespace) Get(ctx context.Context) (*corev1.Namespace, error) {
	opts := metav1.GetOptions{}
	return r.client.CoreV1().Namespaces().Get(ctx, ns, opts)
}

func (r *namespace) Create(ctx context.Context) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}
	opts := metav1.CreateOptions{}
	return r.client.CoreV1().Namespaces().Create(ctx, ns, opts)
}
