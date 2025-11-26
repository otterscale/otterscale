package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
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

func (r *nodeRepo) InternalIP(ctx context.Context, scope string) (string, error) {
	selector := "node-role.kubernetes.io/control-plane"

	nodes, err := r.List(ctx, scope, selector)
	if err != nil {
		return "", err
	}

	for i := range nodes {
		if !r.isNodeReady(&nodes[i]) {
			continue
		}

		if ip := r.getInternalIP(&nodes[i]); ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("no control plane node with InternalIP found")
}

func (r *nodeRepo) isNodeReady(node *cluster.Node) bool {
	for i := range node.Status.Conditions {
		if node.Status.Conditions[i].Type == v1.NodeReady &&
			node.Status.Conditions[i].Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func (r *nodeRepo) getInternalIP(node *cluster.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}
