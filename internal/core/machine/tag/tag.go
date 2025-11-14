package tag

import (
	"context"

	"golang.org/x/sync/errgroup"
)

const (
	BuiltIn                = "built-in"
	Kubernetes             = "kubernetes"
	KubernetesControlPlane = "kubernetes-control-plane"
	KubernetesWorker       = "kubernetes-worker"
	Ceph                   = "ceph"
	CephMON                = "ceph-mon"
	CephOSD                = "ceph-osd"
)

type Tag struct{}

type TagRepo interface {
	List(ctx context.Context) ([]Tag, error)
	Get(ctx context.Context, name string) (*Tag, error)
	Create(ctx context.Context, name, comment string) (*Tag, error)
	Delete(ctx context.Context, name string) error
	AddMachines(ctx context.Context, name string, machineIDs []string) error
	RemoveMachines(ctx context.Context, name string, machineIDs []string) error
}

type TagUseCase struct {
	tag TagRepo
}

func NewTagUseCase(tag TagRepo) *TagUseCase {
	return &TagUseCase{
		tag: tag,
	}
}

func (uc *TagUseCase) ListTags(ctx context.Context) ([]Tag, error) {
	return uc.tag.List(ctx)
}

func (uc *TagUseCase) GetTag(ctx context.Context, name string) (*Tag, error) {
	return uc.tag.Get(ctx, name)
}

func (uc *TagUseCase) CreateTag(ctx context.Context, name, comment string) (*Tag, error) {
	return uc.tag.Create(ctx, name, comment)
}

func (uc *TagUseCase) DeleteTag(ctx context.Context, name string) error {
	return uc.tag.Delete(ctx, name)
}

func (uc *TagUseCase) AddMachineTags(ctx context.Context, id string, tags []string) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.AddMachines(egctx, tag, []string{id})
		})
	}

	return eg.Wait()
}

func (uc *TagUseCase) RemoveMachineTags(ctx context.Context, id string, tags []string) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.RemoveMachines(egctx, tag, []string{id})
		})
	}

	return eg.Wait()
}
