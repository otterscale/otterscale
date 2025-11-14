package kubernetes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"

	"sigs.k8s.io/yaml"

	"github.com/otterscale/otterscale/internal/core/application/chart"
)

const repoIndexCacheDuration = time.Hour * 4

type repoIndex struct {
	file      *repo.IndexFile
	lastFetch time.Time
}

type chartRepo struct {
	kubernetes *Kubernetes

	repoIndexCache sync.Map
}

func NewChartRepo(kubernetes *Kubernetes) (chart.ChartRepo, error) {
	return &chartRepo{
		kubernetes: kubernetes,
	}, nil
}

var _ chart.ChartRepo = (*chartRepo)(nil)

func (r *chartRepo) List(ctx context.Context, url string, useCache bool) ([]chart.Chart, error) {
	if useCache {
		if repoIndex, ok := r.getRepoIndexCache(url); ok {
			return r.buildChartsFromIndex(repoIndex.file), nil
		}
	}

	indexFile, err := r.fetchRepoIndex(ctx, url)
	if err != nil {
		return nil, err
	}

	r.setRepoIndexCache(url, indexFile)

	return r.buildChartsFromIndex(indexFile), nil
}

func (r *chartRepo) Show(chartRef string, format action.ShowOutputFormat) (string, error) {
	client := action.NewShow(format)
	client.SetRegistryClient(r.kubernetes.registryClient)

	chartPath, err := client.LocateChart(chartRef, r.kubernetes.envSettings)
	if err != nil {
		return "", err
	}

	return client.Run(chartPath)
}

func (r *chartRepo) Push(chartRef, remoteOCI string) (string, error) {
	config := &action.Configuration{
		RegistryClient: r.kubernetes.registryClient,
	}

	client := action.NewPushWithOpts(action.WithPushConfig(config))
	client.Settings = r.kubernetes.envSettings

	return client.Run(chartRef, remoteOCI)
}

func (r *chartRepo) Index(dir, url string) error {
	out := filepath.Join(dir, "index.yaml")

	i, err := repo.IndexDirectory(dir, url)
	if err != nil {
		return err
	}

	i.SortEntries()

	return i.WriteFile(out, 0o644) //nolint:mnd // default file permission
}

func (r *chartRepo) GetStableVersion(ctx context.Context, url, name string, useCache bool) (*chart.Version, error) {
	if useCache {
		if repoIndex, ok := r.getRepoIndexCache(url); ok {
			return repoIndex.file.Get(name, "")
		}
	}

	indexFile, err := r.fetchRepoIndex(ctx, url)
	if err != nil {
		return nil, err
	}

	r.setRepoIndexCache(url, indexFile)

	return indexFile.Get(name, "")
}

func (r *chartRepo) buildChartsFromIndex(indexFile *repo.IndexFile) []chart.Chart {
	charts := make([]chart.Chart, 0, len(indexFile.Entries))

	for name, versions := range indexFile.Entries {
		charts = append(charts, chart.Chart{
			Name:     name,
			Versions: versions,
		})
	}

	return charts
}

func (r *chartRepo) getRepoIndexCache(url string) (*repoIndex, bool) {
	v, ok := r.repoIndexCache.Load(url)
	if !ok {
		return nil, false
	}

	repoIndex := v.(*repoIndex)
	if time.Since(repoIndex.lastFetch) >= repoIndexCacheDuration {
		return nil, false
	}

	return repoIndex, true
}

func (r *chartRepo) setRepoIndexCache(url string, file *repo.IndexFile) {
	file.SortEntries()

	r.repoIndexCache.Store(url, &repoIndex{
		file:      file,
		lastFetch: time.Now(),
	})
}

func (r *chartRepo) fetchRepoIndex(ctx context.Context, url string) (*repo.IndexFile, error) {
	if strings.HasPrefix(url, "http") {
		return r.fetchRepoIndexFromWeb(ctx, url)
	}
	return r.fetchRepoIndexFromFile(url)
}

func (r *chartRepo) fetchRepoIndexFromWeb(ctx context.Context, repoURL string) (*repo.IndexFile, error) {
	queryURL, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return nil, err
	}

	queryURL = queryURL.JoinPath("index.yaml")

	data, err := r.httpGet(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	index := new(repo.IndexFile)

	if err := yaml.Unmarshal(data, index); err != nil {
		return nil, err
	}

	return index, nil
}

func (r *chartRepo) fetchRepoIndexFromFile(path string) (*repo.IndexFile, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	path = filepath.Join(path, "index.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	index := new(repo.IndexFile)

	if err := yaml.Unmarshal(data, index); err != nil {
		return nil, err
	}

	return index, nil
}

func (r *chartRepo) httpGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("get %q failed: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get %q failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("get %q failed with code %d", url, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
