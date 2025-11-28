package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

    "github.com/regclient/regclient"
    "github.com/regclient/regclient/config"
	"github.com/regclient/regclient/types/ref"

	"github.com/otterscale/otterscale/internal/core/containerimage"

)

type containerImageRepo struct {
	kubernetes *Kubernetes
}

func NewContainerImageRepo(kubernetes *Kubernetes) (containerimage.ContainerImageRepo, error) {
	return &containerImageRepo{
		kubernetes: kubernetes,
	}, nil
}

var _ containerimage.ContainerImageRepo = (*containerImageRepo)(nil)

func (r *containerImageRepo) newRegClient(endpoint string) *regclient.RegClient {
	exHostLocal := config.Host{
		Name: endpoint,
		TLS:  config.TLSDisabled, 
	}

	return regclient.New(
		regclient.WithConfigHosts([]config.Host{exHostLocal}),
	)
}

func (r *containerImageRepo) ListRepos (ctx context.Context, scope, namespace, endpoint string) ([]string, error){
	rc := r.newRegClient(endpoint)
	defer rc.Close(ctx, ref.Ref{})

	repos := []string{}
	rl, err := rc.RepoList(ctx, endpoint)
	if err != nil {
        return nil, err
    }
	for _, repoName := range rl.Repositories {
		if namespace == "" {
       		repos = append(repos, repoName)
			continue
		}
		prefix := namespace + "/"
		if strings.HasPrefix(repoName, prefix) {
			name := strings.TrimPrefix(repoName, prefix)
			repos = append(repos, name)
		}
    }

	return repos, nil
}

func (r *containerImageRepo) ListTags (ctx context.Context, scope, namespace, endpoint, name string) ([]string, error) {
	rc := r.newRegClient(endpoint)
	defer rc.Close(ctx, ref.Ref{})

	repoPath := name
	if namespace != "" {
		repoPath = fmt.Sprintf("%s/%s", namespace, name)
	}

	imageRef := fmt.Sprintf("%s/%s", endpoint, repoPath)
	refRepo, err := ref.New(imageRef)
	if err != nil {
		return nil, err
	}

	tl, err := rc.TagList(ctx, refRepo)
	if err != nil {
		return nil, err
	}

	return tl.Tags, nil
}

func (r *containerImageRepo) GetSize (ctx context.Context, scope, namespace, endpoint, name, tag string) (int64, error) {
	rc := r.newRegClient(endpoint)
	defer rc.Close(ctx, ref.Ref{})

	repoPath := name
	if namespace != "" {
		repoPath = fmt.Sprintf("%s/%s", namespace, name)
	}

	imageRef := fmt.Sprintf("%s/%s:%s", endpoint, repoPath, tag)
	refImg, err := ref.New(imageRef)
	if err != nil {
		return 0, err
	}

	m, err := rc.ManifestGet(ctx, refImg)
	if err != nil {
		return 0, err
	}

	var totalSize int64

	layers, err := m.GetLayers()
	if err == nil {
		for _, l := range layers {
			totalSize += l.Size
		}
	}

	cfgDesc, err := m.GetConfig()
	if err == nil {
		totalSize += cfgDesc.Size
	}

	return totalSize, nil
}

func (r *containerImageRepo) GetCreateAt (ctx context.Context, scope, namespace, endpoint, name, tag string) (time.Time, error) {
	rc := r.newRegClient(endpoint)
	defer rc.Close(ctx, ref.Ref{})

	repoPath := name
	if namespace != "" {
		repoPath = fmt.Sprintf("%s/%s", namespace, name)
	}

	imageRef := fmt.Sprintf("%s/%s:%s", endpoint, repoPath, tag)
	refImg, err := ref.New(imageRef)
	if err != nil {
		return time.Time{}, err
	}

	cfg, err := rc.ImageConfig(ctx, refImg)
    if err != nil {
        return time.Time{}, err
    }
    if cfg == nil {
        return time.Time{}, nil
    }

    imgCfg := cfg.GetConfig()
    if imgCfg.Created == nil {
        return time.Time{}, nil
    }

	return *imgCfg.Created, nil
}

func (r *containerImageRepo) UploadImage(ctx context.Context, scope, namespace, endpoint, name, tag string, imageTar []byte) error {
    rc := r.newRegClient(endpoint)
    defer rc.Close(ctx, ref.Ref{})

    repoPath := name
    if namespace != "" {
        repoPath = fmt.Sprintf("%s/%s", namespace, name)
    }

    imageRef := fmt.Sprintf("%s/%s:%s", endpoint, repoPath, tag)

    rRef, err := ref.New(imageRef)
    if err != nil {
        return fmt.Errorf("build image ref %q: %w", imageRef, err)
    }

    rs := bytes.NewReader(imageTar)

    if err := rc.ImageImport(ctx, rRef, rs); err != nil {
        return fmt.Errorf("image import to %q failed: %w", imageRef, err)
    }

    return nil
}

func (r *containerImageRepo) DeleteImage (ctx context.Context, scope, namespace, endpoint, name, tag string) error {
	rc := r.newRegClient(endpoint)
    defer rc.Close(ctx, ref.Ref{})

	repoPath := name
    if namespace != "" {
        repoPath = fmt.Sprintf("%s/%s", namespace, name)
    }

	imageRefTag := fmt.Sprintf("%s/%s:%s", endpoint, repoPath, tag)

	rRefTag, err := ref.New(imageRefTag)
    if err != nil {
        return fmt.Errorf("build tag ref %q: %w", imageRefTag, err)
    }

	desc, err := rc.ManifestHead(ctx, rRefTag)
    if err != nil {
        return fmt.Errorf("manifest head for %q failed: %w", imageRefTag, err)
    }

    d := desc.GetDescriptor().Digest
    if d == "" {
        return fmt.Errorf("no digest returned for %q", imageRefTag)
    }

	imageRefDigest := fmt.Sprintf("%s/%s@%s", endpoint, repoPath, d.String())

    rRefDigest, err := ref.New(imageRefDigest)
    if err != nil {
        return fmt.Errorf("build digest ref %q: %w", imageRefDigest, err)
    }

	if err := rc.ManifestDelete(ctx, rRefDigest); err != nil {
        return fmt.Errorf("image delete failed: %w", err)
    }

	return nil
}
