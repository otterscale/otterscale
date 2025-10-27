package kube

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"

	"sigs.k8s.io/yaml"

	oscore "github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/wrap"
)

type helmRepo struct {
	charts    []oscore.Chart
	lastFetch time.Time
}

type helmChart struct {
	kube           *Kube
	repoIndexCache sync.Map
}

func NewHelmChart(kube *Kube) (oscore.ChartRepo, error) {
	return &helmChart{
		kube: kube,
	}, nil
}

var _ oscore.ChartRepo = (*helmChart)(nil)

func (r *helmChart) List(ctx context.Context, url string) ([]oscore.Chart, error) {
	if charts, ok := r.getCachedCharts(url); ok {
		return charts, nil
	}

	indexFile, err := r.fetchRepoIndex(ctx, url)
	if err != nil {
		return nil, err
	}

	charts := r.buildChartsFromIndex(indexFile)
	r.cacheCharts(url, charts)

	return charts, nil
}

func (r *helmChart) Show(chartRef string, format action.ShowOutputFormat) (string, error) {
	client := action.NewShow(format)
	client.SetRegistryClient(r.kube.registryClient)

	chartPath, err := client.LocateChart(chartRef, r.kube.envSettings)
	if err != nil {
		return "", err
	}
	return client.Run(chartPath)
}

func (r *helmChart) Push(chartRef, remoteOCI string) (string, error) {
	config := &action.Configuration{
		RegistryClient: r.kube.registryClient,
	}
	client := action.NewPushWithOpts(action.WithPushConfig(config))
	client.Settings = r.kube.envSettings
	return client.Run(chartRef, remoteOCI)
}

func (r *helmChart) Index(dir, url string) error {
	out := filepath.Join(dir, "index.yaml")

	i, err := repo.IndexDirectory(dir, url)
	if err != nil {
		return err
	}
	i.SortEntries()
	return i.WriteFile(out, 0o644)
}

func (r *helmChart) getCachedCharts(url string) ([]oscore.Chart, bool) {
	v, ok := r.repoIndexCache.Load(url)
	if !ok {
		return nil, false
	}

	helmRepo := v.(*helmRepo)
	if time.Since(helmRepo.lastFetch) >= time.Hour*4 {
		return nil, false
	}

	return helmRepo.charts, true
}

func (r *helmChart) buildChartsFromIndex(indexFile *repo.IndexFile) []oscore.Chart {
	indexFile.SortEntries()
	charts := make([]oscore.Chart, 0, len(indexFile.Entries))
	for name, versions := range indexFile.Entries {
		charts = append(charts, oscore.Chart{
			Name:     name,
			Versions: versions,
		})
	}
	return charts
}

func (r *helmChart) cacheCharts(url string, charts []oscore.Chart) {
	r.repoIndexCache.Store(url, &helmRepo{
		charts:    charts,
		lastFetch: time.Now(),
	})
}

func (r *helmChart) fetchRepoIndex(ctx context.Context, url string) (*repo.IndexFile, error) {
	if strings.HasPrefix(url, "http") {
		return r.fetchRepoIndexFromWeb(ctx, url)
	}
	return r.fetchRepoIndexFromFile(url)
}

func (r *helmChart) fetchRepoIndexFromWeb(ctx context.Context, repoURL string) (*repo.IndexFile, error) {
	queryURL, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return nil, err
	}
	queryURL = queryURL.JoinPath("index.yaml")

	data, err := wrap.HTTPGet(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	f := new(repo.IndexFile)
	if err := yaml.Unmarshal(data, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (r *helmChart) fetchRepoIndexFromFile(path string) (*repo.IndexFile, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, "index.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f := new(repo.IndexFile)
	if err := yaml.Unmarshal(data, f); err != nil {
		return nil, err
	}
	return f, nil
}
