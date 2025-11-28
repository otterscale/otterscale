package orchestrator

import (
	"context"
	"maps"
	"strings"

	"github.com/otterscale/otterscale/internal/core/application/cluster"
)

const domainLabel = "otterscale.com"

type UseCase struct {
	node cluster.NodeRepo
}

func NewUseCase(node cluster.NodeRepo) *UseCase {
	return &UseCase{
		node: node,
	}
}

func (uc *UseCase) ListKubernetesNodeLabels(ctx context.Context, scope, hostname string, all bool) (map[string]string, error) {
	node, err := uc.node.Get(ctx, scope, hostname)
	if err != nil {
		return nil, err
	}

	if !all {
		maps.DeleteFunc(node.Labels, func(k, _ string) bool {
			parts := strings.Split(k, "/")

			return len(parts) < 2 || !strings.HasSuffix(parts[0], domainLabel)
		})
	}

	return node.Labels, nil
}

func (uc *UseCase) UpdateKubernetesNodeLabels(ctx context.Context, scope, hostname string, labels map[string]string) (map[string]string, error) {
	node, err := uc.node.Get(ctx, scope, hostname)
	if err != nil {
		return nil, err
	}

	if node.Labels == nil {
		node.Labels = map[string]string{}
	}

	for k, v := range labels {
		if v == "" {
			delete(node.Labels, k)
		} else {
			node.Labels[k] = v
		}
	}

	node, err = uc.node.Update(ctx, scope, node)
	if err != nil {
		return nil, err
	}

	return node.Labels, nil
}
