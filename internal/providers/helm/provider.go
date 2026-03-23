// Package helm implements the core.HelmRepo interface using the Helm
// v4 Go SDK. It fetches chart metadata (values.yaml, README.md)
// directly from remote HTTP/OCI chart repositories without requiring
// a Kubernetes cluster connection.
package helm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/common"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	chartutil "helm.sh/helm/v4/pkg/chart/v2/util"
	"helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/registry"

	"github.com/otterscale/otterscale/internal/core"
)

var readmeFileNames = []string{"readme.md", "readme.txt", "readme"}

// Repo implements core.HelmRepo using the Helm v4 Go SDK.
// It holds a reusable registry client with auth caching enabled.
type Repo struct {
	registryClient *registry.Client
}

// NewRepo returns a new Helm repository adapter.
func NewRepo() (core.HelmRepo, error) {
	rc, err := registry.NewClient(registry.ClientOptEnableCache(true))
	if err != nil {
		return nil, fmt.Errorf("failed to create helm registry client: %w", err)
	}
	return &Repo{registryClient: rc}, nil
}

// ShowChart fetches a chart from the given repo URL, then extracts
// the values.yaml and README.md content. The chart is loaded once
// and both outputs are extracted from the in-memory structure.
func (r *Repo) ShowChart(_ context.Context, repoURL, chartName, version string) (values, readme []byte, err error) {
	tmpDir, tmpErr := os.MkdirTemp("", "otterscale-helm-")
	if tmpErr != nil {
		return nil, nil, &core.DomainError{
			Code:    core.ErrorCodeInternal,
			Message: "failed to create temporary directory for helm",
			Cause:   tmpErr,
		}
	}
	defer os.RemoveAll(tmpDir)

	settings := cli.New()
	settings.RepositoryConfig = filepath.Join(tmpDir, "repositories.yaml")
	settings.RepositoryCache = filepath.Join(tmpDir, "repository")

	cfg := action.NewConfiguration()
	cfg.RegistryClient = r.registryClient

	show := action.NewShow(action.ShowAll, cfg)
	show.Version = version

	var chartRef string
	if strings.HasPrefix(repoURL, "oci://") {
		chartRef = strings.TrimSuffix(repoURL, "/") + "/" + chartName
	} else {
		show.RepoURL = repoURL
		chartRef = chartName
	}

	chartPath, err := show.LocateChart(chartRef, settings)
	if err != nil {
		return nil, nil, &core.DomainError{
			Code:    classifyHelmError(err),
			Message: fmt.Sprintf("failed to locate chart %s/%s", repoURL, chartName),
			Cause:   err,
		}
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, nil, &core.DomainError{
			Code:    core.ErrorCodeInternal,
			Message: fmt.Sprintf("failed to load chart %s/%s", repoURL, chartName),
			Cause:   err,
		}
	}

	for _, f := range chart.Raw {
		if f.Name == chartutil.ValuesfileName {
			values = f.Data
			break
		}
	}

	if f := findReadme(chart.Files); f != nil {
		readme = f.Data
	}

	return values, readme, nil
}

func findReadme(files []*common.File) *common.File {
	for _, file := range files {
		if file == nil {
			continue
		}
		for _, n := range readmeFileNames {
			if strings.EqualFold(file.Name, n) {
				return file
			}
		}
	}
	return nil
}

// classifyHelmError maps a Helm SDK error to a domain error code.
func classifyHelmError(err error) core.ErrorCode {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "not found"):
		return core.ErrorCodeNotFound
	case strings.Contains(msg, "invalid"):
		return core.ErrorCodeInvalidArgument
	default:
		return core.ErrorCodeInternal
	}
}
