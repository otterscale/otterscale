package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/cluster"
)

type nodeRepo struct {
	kubernetes *Kubernetes
}

func NewNodeRepo(kubernetes *Kubernetes) cluster.NodeRepo {
	return &nodeRepo{
		kubernetes: kubernetes,
	}
}

var _ cluster.NodeRepo = (*nodeRepo)(nil)

func (r *nodeRepo) List(ctx context.Context, scope, selector string) ([]cluster.Node, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CoreV1().Nodes().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *nodeRepo) Get(ctx context.Context, scope, name string) (*cluster.Node, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CoreV1().Nodes().Get(ctx, name, opts)
}

func (r *nodeRepo) Update(ctx context.Context, scope string, n *cluster.Node) (*cluster.Node, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.CoreV1().Nodes().Update(ctx, n, opts)
}
