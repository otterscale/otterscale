package registry

import (
	"context"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/otterscale/otterscale/internal/core/registry"
)

type repositoryRepo struct {
	registry *Registry
}

func NewRepositoryRepo(registry *Registry) registry.RepositoryRepo {
	return &repositoryRepo{
		registry: registry,
	}
}

var _ registry.RepositoryRepo = (*repositoryRepo)(nil)

// FIXME: implement pagination
func (r *repositoryRepo) List(ctx context.Context, scope string) ([]registry.Repository, error) {
	client, err := r.registry.client(scope)
	if err != nil {
		return nil, err
	}

	repositories := []registry.Repository{}

	err = client.Repositories(ctx, "", func(repos []string) error {
		for _, repo := range repos {
			repository, err := r.buildRepository(ctx, client, repo)
			if err != nil {
				return err
			}

			repositories = append(repositories, *repository)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func (r *repositoryRepo) GetRegistryURL(scope string) (string, error) {
	client, err := r.registry.client(scope)
	if err != nil {
		return "", err
	}

	return client.Reference.Registry, nil
}

func (r *repositoryRepo) buildRepository(ctx context.Context, client *remote.Registry, repository string) (*registry.Repository, error) {
	repo, err := client.Repository(ctx, repository)
	if err != nil {
		return nil, err
	}

	manifestCount := uint32(0)
	sizeBytes := uint64(0)
	latestTag := ""

	err = repo.Tags(ctx, "", func(tags []string) error {
		var tagVersions []*semver.Version

		for _, tag := range tags {
			manifestCount++

			reference := client.Reference.Registry + "/" + repository + ":" + tag

			_, ociManifest, err := fetchManifest(ctx, repo, reference)
			if err != nil {
				return err
			}

			sizeBytes += calculateLayerSize(ociManifest)

			// Change underscore (_) back to plus (+) for Helm
			// See https://github.com/helm/helm/issues/10166
			tagVersion, err := semver.StrictNewVersion(strings.ReplaceAll(tag, "_", "+"))
			if err == nil {
				tagVersions = append(tagVersions, tagVersion)
			}
		}

		sort.Sort(sort.Reverse(semver.Collection(tagVersions)))

		if len(tagVersions) > 0 {
			latestTag = tagVersions[0].Original()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &registry.Repository{
		Name:          repository,
		ManifestCount: manifestCount,
		SizeBytes:     sizeBytes,
		LatestTag:     latestTag,
	}, nil
}
