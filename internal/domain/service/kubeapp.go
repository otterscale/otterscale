package service

import (
	"context"

	"golang.org/x/sync/errgroup"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *KubeService) ListApplications(ctx context.Context, uuid, cluster string) (*model.Applications, error) {
	if err := s.ensureClient(ctx, uuid, cluster); err != nil {
		return nil, err
	}
	var (
		deployments            *appv1.DeploymentList
		services               *corev1.ServiceList
		pods                   *corev1.PodList
		persistentVolumeClaims *corev1.PersistentVolumeClaimList
	)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		dpl, err := s.deployment.List(ctx, cluster, "")
		if err == nil {
			deployments = dpl
		}
		return err
	})
	eg.Go(func() error {
		svc, err := s.service.List(ctx, cluster, "")
		if err == nil {
			services = svc
		}
		return err
	})
	eg.Go(func() error {
		pod, err := s.pod.List(ctx, cluster, "")
		if err == nil {
			pods = pod
		}
		return err
	})
	eg.Go(func() error {
		pvc, err := s.persistentVolumeClaim.List(ctx, cluster, "")
		if err == nil {
			persistentVolumeClaims = pvc
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &model.Applications{
		Deployments:            deployments,
		Services:               services,
		Pods:                   pods,
		PersistentVolumeClaims: persistentVolumeClaims,
	}, nil
}
