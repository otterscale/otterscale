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

	oscore "github.com/openhdc/otterscale/internal/core"
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

func NewHelmChart(kube *Kube) (oscore.ChartRepo, error) {
	return &helmChart{
		kube: kube,
	}, nil
}

var _ oscore.ChartRepo = (*helmChart)(nil)

func (r *helmChart) List(ctx context.Context) ([]oscore.Chart, error) {
	urls := r.kube.helmRepoURLs()
	eg, ctx := errgroup.WithContext(ctx)
	result := make([]*repo.IndexFile, len(urls))
	for i := range urls {
		eg.Go(func() error {
			v, ok := r.repoIndexCache.Load(urls[i])
			if ok {
				helmRepo := v.(*helmRepo)
				if time.Since(helmRepo.lastFetch) < time.Hour*2 {
					result[i] = helmRepo.indexFile
					return nil
				}
			}
			indexFile, err := r.fetchRepoIndex(ctx, urls[i])
			if err == nil {
				indexFile.SortEntries()
				r.repoIndexCache.Store(urls[i], &helmRepo{
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

	charts := []oscore.Chart{}
	for i := range result {
		for name := range result[i].Entries {
			charts = append(charts, oscore.Chart{
				Name:     name,
				Versions: result[i].Entries[name],
			})
		}
	}
	return charts, nil
}

func (r *helmChart) Show(chartRef string, format action.ShowOutputFormat) (string, error) {
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
	return r.fetchRepoIndexFromFile(url)
}

func (r *helmChart) fetchRepoIndexFromWeb(ctx context.Context, repoURL string) (*repo.IndexFile, error) {
	queryURL, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return nil, err
	}
	queryURL = queryURL.JoinPath("index.yaml")

	data, err := utils.HTTPGet(ctx, queryURL.String())
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
