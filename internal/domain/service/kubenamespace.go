package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

type KubeNamespace interface {
	Get(ctx context.Context, cluster, name string) (*corev1.Namespace, error)
	Create(ctx context.Context, cluster, name string) (*corev1.Namespace, error)
}

func (s *KubeService) CreateNamespace(ctx context.Context, cluster, name string) error {
	_, err := s.namespace.Create(ctx, cluster, name)
	if err != nil {
		return err
	}
	return nil
}
