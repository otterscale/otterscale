package registry

import (
	"context"

	"github.com/otterscale/otterscale/internal/core/registry/chart"
	"github.com/otterscale/otterscale/internal/core/registry/img"
)

type Manifest struct {
	Repository string
	Tag        string
	Digest     string
	SizeBytes  uint64
	Image      *img.Image
	Chart      *chart.Chart
}

type ManifestRepo interface {
	List(ctx context.Context, scope, repository string) ([]Manifest, error)
	Get(ctx context.Context, scope, repository, tag string) (*Manifest, error)
	Delete(ctx context.Context, scope, repository, digest string) error
}
