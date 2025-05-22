package kube

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"

	"sigs.k8s.io/yaml"

	"github.com/openhdc/otterscale/internal/domain/service"
	"github.com/openhdc/otterscale/internal/utils"
)

type helmRepo struct {
	indexFile *repo.IndexFile
	lastFetch time.Time
}

type helmChart struct {
	kube           *Kube
	repoIndexCache sync.Map
}

func NewHelmChart(kube *Kube) (service.KubeHelmChart, error) {
	return &helmChart{
		kube: kube,
	}, nil
}

var _ service.KubeHelm = (*helm)(nil)

func (r *helmChart) ListCharts(ctx context.Context) ([]*repo.IndexFile, error) {
	urls := r.kube.helmRepoURLs()
	eg, ctx := errgroup.WithContext(ctx)
	result := make([]*repo.IndexFile, len(urls))
	for i := range urls {
		url := urls[i]
		eg.Go(func() error {
			v, ok := r.repoIndexCache.Load(url)
			if ok {
				helmRepo := v.(*helmRepo)
				if time.Since(helmRepo.lastFetch) < time.Hour*2 {
					result[i] = helmRepo.indexFile
					return nil
				}
			}
			indexFile, err := r.fetchRepoIndex(ctx, url)
			if err == nil {
				indexFile.SortEntries()
				r.repoIndexCache.Store(url, &helmRepo{
					indexFile: indexFile,
					lastFetch: time.Now(),
				})
				result[i] = indexFile
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *helmChart) GetChart(chartRef string, format action.ShowOutputFormat) (string, error) {
	client := action.NewShow(format)
	client.SetRegistryClient(r.kube.registryClient)

	chartPath, err := client.ChartPathOptions.LocateChart(chartRef, r.kube.envSettings)
	if err != nil {
		return "", err
	}
	return client.Run(chartPath)
}

func (r *helmChart) fetchRepoIndex(ctx context.Context, url string) (*repo.IndexFile, error) {
	if strings.HasPrefix(url, "http") {
		return r.fetchRepoIndexFromWeb(ctx, url)
	}
	return r.fetchRepoIndexFromFile(ctx, url)
}

func (r *helmChart) fetchRepoIndexFromWeb(ctx context.Context, repoURL string) (*repo.IndexFile, error) {
	queryURL, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return nil, err
	}
	queryURL = queryURL.JoinPath("index.yaml")

	data, err := utils.Get(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	f := new(repo.IndexFile)
	if err := yaml.Unmarshal(data, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (r *helmChart) fetchRepoIndexFromFile(ctx context.Context, path string) (*repo.IndexFile, error) {
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
