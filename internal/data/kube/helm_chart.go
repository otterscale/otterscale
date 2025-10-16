package kube

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"

	"sigs.k8s.io/yaml"

	oscore "github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/wrap"
)

const (
	localOCIChartsDir = "./charts"
	defaultDirPerm    = 0o755
	defaultFilePerm   = 0o600
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
	// Get configured repository URLs
	configuredURLs := r.kube.helmRepoURLs()

	// Add local OCI upload directory as an additional source
	allURLs := make([]string, len(configuredURLs)+1)
	copy(allURLs, configuredURLs)
	allURLs[len(configuredURLs)] = localOCIChartsDir

	eg, egctx := errgroup.WithContext(ctx)
	result := make([]*repo.IndexFile, len(allURLs))

	for i := range allURLs {
		eg.Go(func() error {
			v, ok := r.repoIndexCache.Load(allURLs[i])
			if ok {
				helmRepo := v.(*helmRepo)
				cacheDuration := time.Hour * 24
				if allURLs[i] == localOCIChartsDir {
					cacheDuration = time.Minute * 1
				}
				if time.Since(helmRepo.lastFetch) < cacheDuration {
					result[i] = helmRepo.indexFile
					return nil
				}
			}
			indexFile, err := r.fetchRepoIndex(egctx, allURLs[i])
			if err == nil {
				indexFile.SortEntries()
				r.repoIndexCache.Store(allURLs[i], &helmRepo{
					indexFile: indexFile,
					lastFetch: time.Now(),
				})
				result[i] = indexFile
			} else if allURLs[i] == localOCIChartsDir {
				return nil
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	charts := []oscore.Chart{}
	for i := range result {
		if result[i] != nil { // Check for nil in case local OCI directory was skipped
			for name := range result[i].Entries {
				charts = append(charts, oscore.Chart{
					Name:     name,
					Versions: result[i].Entries[name],
				})
			}
		}
	}
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

// UploadChart uploads a Helm chart to OCI registry and generates/merges index.yaml
func (r *helmChart) UploadChart(ctx context.Context, ociRegistryURL, chartName, chartVersion string, chartContent []byte) error {
	helmChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:        chartName,
			Version:     chartVersion,
			Description: "Chart uploaded to OCI registry",
		},
	}

	chartRef, err := r.pushToOCI(ctx, ociRegistryURL, chartName, chartVersion, chartContent)
	if err != nil {
		return fmt.Errorf("failed to push chart to OCI registry: %w", err)
	}

	err = r.generateAndMergeIndexYAML(helmChart, chartRef, chartContent)
	if err != nil {
		return fmt.Errorf("failed to generate index.yaml: %w", err)
	}

	return nil
}

// pushToOCI pushes the chart to OCI registry using Helm's OCI support
func (r *helmChart) pushToOCI(ctx context.Context, ociRegistryURL, chartName, chartVersion string, chartContent []byte) (string, error) {
	registryURL := ociRegistryURL
	if !strings.HasPrefix(registryURL, "oci://") {
		registryURL = "oci://" + strings.TrimPrefix(registryURL, "oci://")
	}

	// Clean chart name and version to avoid invalid characters
	cleanChartName := strings.ToLower(strings.ReplaceAll(chartName, "_", "-"))
	cleanChartVersion := strings.ReplaceAll(chartVersion, "+", "-")

	chartRef := fmt.Sprintf("%s/%s:%s", registryURL, cleanChartName, cleanChartVersion)

	// Create a temporary file for the chart
	tmpFile, err := os.CreateTemp("", "helm-chart-*.tgz")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(chartContent); err != nil {
		return "", fmt.Errorf("failed to write chart to temporary file: %w", err)
	}
	tmpFile.Close()

	// Create action configuration with registry client
	cfg := &action.Configuration{
		RegistryClient: r.kube.registryClient,
	}

	pushAction := action.NewPushWithOpts(action.WithPushConfig(cfg))
	pushAction.Settings = r.kube.envSettings

	lastSlash := strings.LastIndex(chartRef, "/")
	if lastSlash == -1 {
		return "", fmt.Errorf("invalid chart reference format: %s", chartRef)
	}

	registryURL = chartRef[:lastSlash]

	_, err = pushAction.Run(tmpFile.Name(), registryURL)
	if err != nil {
		return "", fmt.Errorf("failed to push chart %s: %w", chartRef, err)
	}

	return chartRef, nil
}

func (r *helmChart) generateAndMergeIndexYAML(chart *chart.Chart, chartRef string, chartContent []byte) error {
	indexDir := localOCIChartsDir
	indexPath := filepath.Join(indexDir, "index.yaml")

	if err := os.MkdirAll(indexDir, defaultDirPerm); err != nil {
		return fmt.Errorf("failed to create charts directory: %w", err)
	}

	// Create a new index from the current chart directory using Helm's official method
	newIndex, err := repo.IndexDirectory(indexDir, "")
	if err != nil {
		newIndex = repo.NewIndexFile()
	}

	// Add current chart entry to the new index
	chartVersion := &repo.ChartVersion{
		Metadata: chart.Metadata,
		URLs:     []string{chartRef},
		Created:  time.Now(),
		Digest:   fmt.Sprintf("sha256:%x", sha256.Sum256(chartContent)),
	}

	chartName := chart.Metadata.Name
	if newIndex.Entries == nil {
		newIndex.Entries = make(map[string]repo.ChartVersions)
	}
	if newIndex.Entries[chartName] == nil {
		newIndex.Entries[chartName] = repo.ChartVersions{}
	}

	// Check if this version already exists and update or add
	var found bool
	for i, existing := range newIndex.Entries[chartName] {
		if existing.Version == chart.Metadata.Version {
			newIndex.Entries[chartName][i] = chartVersion
			found = true
			break
		}
	}
	if !found {
		newIndex.Entries[chartName] = append(newIndex.Entries[chartName], chartVersion)
	}

	// Merge with existing index if it exists (using Helm's official merge method)
	if _, err := os.Stat(indexPath); err == nil {
		existingIndex, err := repo.LoadIndexFile(indexPath)
		if err != nil {
			return fmt.Errorf("failed to load existing index file: %w", err)
		}
		newIndex.Merge(existingIndex)
	}

	newIndex.SortEntries()
	newIndex.Generated = time.Now()

	return newIndex.WriteFile(indexPath, defaultFilePerm)
}
