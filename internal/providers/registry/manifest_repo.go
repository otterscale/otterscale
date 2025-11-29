package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	helmregistry "helm.sh/helm/v3/pkg/registry"
	orasregistry "oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/otterscale/otterscale/internal/core/registry"
)

type manifestRepo struct {
	registry *Registry
}

func NewManifestRepo(registry *Registry) registry.ManifestRepo {
	return &manifestRepo{registry: registry}
}

var _ registry.ManifestRepo = (*manifestRepo)(nil)

func (r *manifestRepo) List(ctx context.Context, scope, repository string) ([]registry.Manifest, error) {
	client, err := r.registry.client(scope)
	if err != nil {
		return nil, err
	}

	repo, err := client.Repository(ctx, repository)
	if err != nil {
		return nil, err
	}

	manifests := []registry.Manifest{}

	err = repo.Tags(ctx, "", func(tags []string) error {
		for _, tag := range tags {
			manifest, err := r.buildManifest(ctx, client, repo, repository, tag)
			if err != nil {
				return err
			}

			manifests = append(manifests, manifest)
		}
		return nil
	})

	return manifests, err
}

func (r *manifestRepo) Delete(ctx context.Context, scope, repository, digestStr string) error {
	client, err := r.registry.client(scope)
	if err != nil {
		return err
	}

	repo, err := client.Repository(ctx, repository)
	if err != nil {
		return err
	}

	reference := fmt.Sprintf("%s/%s@%s", client.Reference.Registry, repository, digestStr)

	desc, err := repo.Resolve(ctx, reference)
	if err != nil {
		return err
	}

	return repo.Delete(ctx, desc)
}

func (r *manifestRepo) buildManifest(ctx context.Context, client *remote.Registry, repo orasregistry.Repository, repository, tag string) (registry.Manifest, error) {
	reference := fmt.Sprintf("%s/%s:%s", client.Reference.Registry, repository, tag)

	digest, ociManifest, err := fetchManifest(ctx, repo, reference)
	if err != nil {
		return registry.Manifest{}, err
	}

	image, chart, err := fetchConfig(ctx, repo, ociManifest.Config)
	if err != nil {
		return registry.Manifest{}, err
	}

	return registry.Manifest{
		Repository: repository,
		Tag:        tag,
		Digest:     digest,
		SizeBytes:  calculateLayerSize(ociManifest),
		Image:      image,
		Chart:      chart,
	}, nil
}

func fetchManifest(ctx context.Context, repo orasregistry.Repository, reference string) (string, *ocispec.Manifest, error) {
	desc, reader, err := repo.FetchReference(ctx, reference)
	if err != nil {
		return "", nil, err
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", nil, err
	}

	var manifest ocispec.Manifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return "", nil, err
	}

	return desc.Digest.String(), &manifest, nil
}

func fetchConfig(ctx context.Context, repo orasregistry.Repository, config ocispec.Descriptor) (*registry.Image, *registry.Chart, error) {
	reader, err := repo.Fetch(ctx, config)
	if err != nil {
		return nil, nil, err
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}

	return parseConfig(content, config.MediaType)
}

func parseConfig(content []byte, mediaType string) (*registry.Image, *registry.Chart, error) {
	var (
		image registry.Image
		chart registry.Chart
	)

	switch mediaType {
	case helmregistry.ConfigMediaType:
		if err := json.Unmarshal(content, &chart); err != nil {
			return nil, nil, err
		}
		return nil, &chart, nil

	case ocispec.MediaTypeImageConfig, "application/vnd.docker.container.image.v1+json":
		if err := json.Unmarshal(content, &image); err != nil {
			return nil, nil, err
		}
		return &image, nil, nil
	}

	return nil, nil, nil
}

func calculateLayerSize(manifest *ocispec.Manifest) uint64 {
	var totalSize uint64
	for _, layer := range manifest.Layers {
		totalSize += uint64(layer.Size)
	}
	return totalSize
}
