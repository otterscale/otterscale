package orchestrator

import (
	"context"
	"maps"
	"strings"

	"github.com/otterscale/otterscale/internal/core/application/cluster"
)

const DomainLabel = "otterscale.com"

type OrchestratorUseCase struct {
	node cluster.NodeRepo
}

func NewOrchestratorUseCase(node cluster.NodeRepo) *OrchestratorUseCase {
	return &OrchestratorUseCase{
		node: node,
	}
}

func (uc *OrchestratorUseCase) ListKubernetesNodeLabels(ctx context.Context, scope, hostname string, all bool) (map[string]string, error) {
	node, err := uc.node.Get(ctx, scope, hostname)
	if err != nil {
		return nil, err
	}

	if !all {
		maps.DeleteFunc(node.Labels, func(k, _ string) bool {
			parts := strings.Split(k, "/")

			return len(parts) < 2 || !strings.HasSuffix(parts[0], DomainLabel)
		})
	}

	return node.Labels, nil
}

func (uc *OrchestratorUseCase) UpdateKubernetesNodeLabels(ctx context.Context, scope, hostname string, labels map[string]string) (map[string]string, error) {
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
