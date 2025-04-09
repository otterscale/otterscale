package kube

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/domain/service"
	"github.com/openhdc/openhdc/internal/env"
)

const searchMaxScore = 25

type helm struct {
	kubeMap   KubeMap
	settings  *cli.EnvSettings
	providers getter.Providers
}

func NewHelm(kubeMap KubeMap) service.KubeHelm {
	settings := cli.New()
	settings.RepositoryConfig = env.GetOrDefault(env.OPENHDC_HELM_REPOSITORY_CONFIG, helmpath.ConfigPath("repositories.yaml"))
	return &helm{
		kubeMap:   kubeMap,
		settings:  settings,
		providers: getter.All(settings),
	}
}

var _ service.KubeHelm = (*helm)(nil)

func (r *helm) ListReleases(ctx context.Context, cluster, namespace string) ([]*release.Release, error) {
	config, err := r.kubeMap.GetHelmConfig(cluster, namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewList(config)
	client.Deployed = true
	return client.Run()
}

func (r *helm) ListRepositories(ctx context.Context) ([]*model.HelmRepo, error) {
	rf, err := repo.LoadFile(r.settings.RepositoryConfig)
	if err != nil {
		return nil, err
	}
	ret := []*model.HelmRepo{}
	for _, re := range rf.Repositories {
		ret = append(ret, &model.HelmRepo{
			Entry:      re,
			ChartNames: r.chartNames(re.Name),
		})
	}
	return ret, nil
}

func (r *helm) UpdateRepositoryCharts(ctx context.Context, name string) (*model.HelmRepo, error) {
	rf, err := repo.LoadFile(r.settings.RepositoryConfig)
	if err != nil {
		return nil, err
	}
	for _, re := range rf.Repositories {
		if re.Name != name {
			continue
		}
		cr, err := repo.NewChartRepository(re, r.providers)
		if err != nil {
			return nil, err
		}
		if _, err := cr.DownloadIndexFile(); err != nil {
			return nil, err
		}
		return &model.HelmRepo{
			Entry:      re,
			ChartNames: r.chartNames(re.Name),
		}, nil
	}
	return nil, fmt.Errorf("helm repo %q not found", name)
}

func (r *helm) ListChartVersions(ctx context.Context) (map[string]repo.ChartVersions, error) {
	rf, err := repo.LoadFile(r.settings.RepositoryConfig)
	if err != nil {
		return nil, err
	}
	ret := map[string]repo.ChartVersions{}
	for _, re := range rf.Repositories {
		path := filepath.Join(r.settings.RepositoryCache, helmpath.CacheIndexFile(re.Name))
		idx, err := repo.LoadIndexFile(path)
		if err != nil {
			continue
		}
		idx.SortEntries()
		maps.Copy(ret, idx.Entries)
	}
	return ret, nil
}

func (r *helm) chartNames(repoName string) []string {
	var charts []string

	path := filepath.Join(r.settings.RepositoryCache, helmpath.CacheChartsFile(repoName))
	content, err := os.ReadFile(path)
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(content))
		for scanner.Scan() {
			fullName := fmt.Sprintf("%s/%s", repoName, scanner.Text())
			charts = append(charts, fullName)
		}
		return charts
	}

	if os.IsNotExist(err) {
		// If there is no cached charts file, fallback to the full index file.
		// This is much slower but can happen after the caching feature is first
		// installed but before the user  does a 'helm repo update' to generate the
		// first cached charts file.
		path = filepath.Join(r.settings.RepositoryCache, helmpath.CacheIndexFile(repoName))
		if indexFile, err := repo.LoadIndexFile(path); err == nil {
			for name := range indexFile.Entries {
				fullName := fmt.Sprintf("%s/%s", repoName, name)
				charts = append(charts, fullName)
			}
			return charts
		}
	}

	return []string{}
}
