package registry

import (
	"context"
)

type Repository struct {
	Name          string
	ManifestCount uint32
	SizeBytes     uint64
	LatestTag     string
}

type RepositoryRepo interface {
	List(ctx context.Context, scope string) ([]Repository, error)
	GetRegistryURL(scope string) (string, error)
}
