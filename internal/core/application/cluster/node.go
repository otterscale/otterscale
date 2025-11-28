package cluster

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

// Node represents a Kubernetes Node resource.
type Node = v1.Node

type NodeRepo interface {
	List(ctx context.Context, scope, selector string) ([]Node, error)
	Get(ctx context.Context, scope, name string) (*Node, error)
	Update(ctx context.Context, scope string, n *Node) (*Node, error)
	InternalIP(ctx context.Context, scope string) (string, error)
}
