package tag

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
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

// Tag represents a MAAS Tag resource.
type Tag = entity.Tag

//nolint:revive // allows this exported interface name for specific domain clarity.
type TagRepo interface {
	List(ctx context.Context) ([]Tag, error)
	Get(ctx context.Context, name string) (*Tag, error)
	Create(ctx context.Context, name, comment string) (*Tag, error)
	Delete(ctx context.Context, name string) error
	AddMachines(ctx context.Context, name string, machineIDs []string) error
	RemoveMachines(ctx context.Context, name string, machineIDs []string) error
}

type UseCase struct {
	tag TagRepo
}

func NewUseCase(tag TagRepo) *UseCase {
	return &UseCase{
		tag: tag,
	}
}

func (uc *UseCase) ListTags(ctx context.Context) ([]Tag, error) {
	return uc.tag.List(ctx)
}

func (uc *UseCase) GetTag(ctx context.Context, name string) (*Tag, error) {
	return uc.tag.Get(ctx, name)
}

func (uc *UseCase) CreateTag(ctx context.Context, name, comment string) (*Tag, error) {
	return uc.tag.Create(ctx, name, comment)
}

func (uc *UseCase) DeleteTag(ctx context.Context, name string) error {
	return uc.tag.Delete(ctx, name)
}

func (uc *UseCase) AddMachineTags(ctx context.Context, id string, tags []string) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.AddMachines(egctx, tag, []string{id})
		})
	}

	return eg.Wait()
}

func (uc *UseCase) RemoveMachineTags(ctx context.Context, id string, tags []string) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.RemoveMachines(egctx, tag, []string{id})
		})
	}

	return eg.Wait()
}
