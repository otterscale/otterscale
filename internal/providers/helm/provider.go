// Package helm implements the core.HelmRepo interface using the Helm
// v4 Go SDK. It fetches chart metadata (values.yaml, README.md)
// directly from remote HTTP/OCI chart repositories without requiring
// a Kubernetes cluster connection.
package helm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/common"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"helm.sh/helm/v4/pkg/chart/v2/util"
	"helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/helmpath"
	"helm.sh/helm/v4/pkg/registry"

	"github.com/otterscale/otterscale/internal/core"
)

var readmeFileNames = []string{"readme.md", "readme.txt", "readme"}

// showChartTimeout bounds a single chart fetch. The Helm SDK offers no
// way to cancel one, so this is what keeps a request from waiting on an
// unresponsive repository when the caller set no deadline of its own.
const showChartTimeout = 60 * time.Second

// Repo implements core.HelmRepo using the Helm v4 Go SDK.
// It holds a reusable registry client with auth caching enabled.
type Repo struct {
	registryClient *registry.Client
	settings       *cli.EnvSettings
}

// NewRepo returns a new Helm repository adapter.
func NewRepo() (core.HelmRepo, error) {
	rc, err := registry.NewClient(registry.ClientOptEnableCache(true))
	if err != nil {
		return nil, fmt.Errorf("failed to create helm registry client: %w", err)
	}

	if err := setHelmHome(); err != nil {
		return nil, err
	}

	// cli.New reads the HELM_* variables set above, so the settings
	// carry the same paths without any later call depending on the
	// environment being read at the right moment.
	return &Repo{registryClient: rc, settings: cli.New()}, nil
}

// helmHome returns the base directory Helm may write to: a stable
// per-user path under the system temp directory. A fixed name is
// deliberate — a fresh directory on every start would leak one cache
// tree per restart, whereas reusing this one also preserves the chart
// cache across restarts.
func helmHome() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("otterscale-helm-%d", os.Getuid()))
}

// setHelmHome points Helm's cache, config and data directories at a
// writable location, without overriding anything the operator has
// already configured.
//
// The environment is the only lever available: parts of the SDK reached
// through LocateChart — repo.NewChartRepository in particular — build
// their cache path from helmpath.CachePath instead of the settings they
// are handed, and helmpath resolves the HELM_*_HOME variables at call
// time. The default resolves under $HOME, which a distroless container
// running as nonroot cannot write to.
func setHelmHome() error {
	base := helmHome()

	for _, dir := range []struct{ env, sub string }{
		{helmpath.CacheHomeEnvVar, "cache"},
		{helmpath.ConfigHomeEnvVar, "config"},
		{helmpath.DataHomeEnvVar, "data"},
	} {
		if os.Getenv(dir.env) != "" {
			continue // operator-provided; leave it alone
		}
		path := filepath.Join(base, dir.sub)
		const dirPerm = 0o700 // this process is the only intended user
		if err := os.MkdirAll(path, dirPerm); err != nil {
			return fmt.Errorf("failed to create helm %s directory: %w", dir.sub, err)
		}
		if err := os.Setenv(dir.env, path); err != nil {
			return fmt.Errorf("failed to set %s: %w", dir.env, err)
		}
	}

	return nil
}

// ShowChart fetches a chart from the given repo URL and extracts its
// values.yaml and README.md.
//
// The Helm SDK call is blocking and takes no context, so it runs on its
// own goroutine and this returns as soon as the caller gives up. The
// fetch itself keeps running until it finishes on its own; its result
// is then discarded.
func (r *Repo) ShowChart(ctx context.Context, repoURL, chartName, version string) (values, readme []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, showChartTimeout)
	defer cancel()

	content, err := awaitWithContext(
		ctx,
		fmt.Sprintf("fetching chart %s/%s", repoURL, chartName),
		func() (chartContent, error) {
			v, rm, err := r.showChart(repoURL, chartName, version)
			return chartContent{values: v, readme: rm}, err
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return content.values, content.readme, nil
}

// chartContent is what a single chart fetch produces.
type chartContent struct {
	values []byte
	readme []byte
}

// awaitWithContext runs fn on its own goroutine and returns as soon as
// ctx is done. It is how a blocking SDK call that takes no context
// becomes cancellable from the caller's side; fn still runs to
// completion, and its result is then discarded. what names the
// operation for the error message.
func awaitWithContext[T any](ctx context.Context, what string, fn func() (T, error)) (T, error) {
	type result struct {
		value T
		err   error
	}

	done := make(chan result, 1)
	go func() {
		value, err := fn()
		done <- result{value: value, err: err}
	}()

	select {
	case res := <-done:
		return res.value, res.err
	case <-ctx.Done():
		var zero T
		return zero, &core.DomainError{
			Code:    ctxErrorCode(ctx.Err()),
			Message: fmt.Sprintf("gave up %s", what),
			Cause:   ctx.Err(),
		}
	}
}

// ctxErrorCode maps a context error to its domain equivalent.
func ctxErrorCode(err error) core.ErrorCode {
	if errors.Is(err, context.DeadlineExceeded) {
		return core.ErrorCodeDeadlineExceeded
	}
	return core.ErrorCodeCanceled
}

// showChart performs the blocking fetch. The chart is loaded once and
// both outputs are extracted from the in-memory structure.
func (r *Repo) showChart(repoURL, chartName, version string) (values, readme []byte, err error) {
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

	chartPath, err := show.LocateChart(chartRef, r.settings)
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
		if f.Name == util.ValuesfileName {
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
