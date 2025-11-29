package registry

import (
	"context"
)

type Manifest struct {
	Repository string
	Tag        string
	Digest     string
	SizeBytes  uint64
	Image      *Image
	Chart      *Chart
}

type ManifestRepo interface {
	List(ctx context.Context, scope, repository string) ([]Manifest, error)
	Delete(ctx context.Context, scope, repository, digest string) error
}
