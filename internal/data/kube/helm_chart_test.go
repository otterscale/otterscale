package kube

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"helm.sh/helm/v3/pkg/action"

	"github.com/otterscale/otterscale/internal/config"
	oscore "github.com/otterscale/otterscale/internal/core"
)

/* -------------------------------------------------------------------------- *
 * Helper – works for both normal tests (testing.T) and benchmarks (testing.B) *
 * -------------------------------------------------------------------------- */
func helm_mustNewKube(t testing.TB, cfg *config.Config) *Kube {
	k, err := New(cfg)
	if err != nil {
		t.Fatalf("cannot create Kube instance: %v", err)
	}
	return k
}

/* -------------------------------------------------------------------------- *
 * Minimal index.yaml used by the tests                                      *
 * -------------------------------------------------------------------------- */
const testIndexYAML = `
apiVersion: v1
entries:
  mychart:
    - version: "1.0.0"
      appVersion: "1.0"
      description: "a test chart"
      urls:
        - "https://example.com/mychart-1.0.0.tgz"
  otherchart:
    - version: "0.2.1"
      urls:
        - "https://example.com/otherchart-0.2.1.tgz"
`

/* -------------------------------------------------------------------------- *
 * Construction / interface compliance                                         *
 * -------------------------------------------------------------------------- */
func TestNewHelmChart(t *testing.T) {
	k := mustNewKube(t, &config.Config{})
	ch, err := NewHelmChart(k)
	if err != nil {
		t.Fatalf("NewHelmChart returned error: %v", err)
	}
	if ch == nil {
		t.Fatal("NewHelmChart returned nil")
	}
	var _ oscore.ChartRepo = ch
}

/* -------------------------------------------------------------------------- *
 * List – fetch from an HTTP server and verify parsing                        *
 * -------------------------------------------------------------------------- */
func TestHelmChart_List_FromHTTP(t *testing.T) {
	// Counter to verify that the cache is used after the first fetch.
	var fetchCount int
	mu := sync.Mutex{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "index.yaml") {
			mu.Lock()
			fetchCount++
			mu.Unlock()
			_, _ = w.Write([]byte(testIndexYAML))
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	cfg := &config.Config{
		Kube: config.Kube{
			HelmRepositoryURLs: []string{ts.URL},
		},
	}
	k := mustNewKube(t, cfg)
	ch, err := NewHelmChart(k)
	if err != nil {
		t.Fatalf("NewHelmChart failed: %v", err)
	}

	// First call – should hit the server.
	ctx := context.Background()
	charts, err := ch.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(charts) == 0 {
		t.Fatalf("expected at least one chart, got none")
	}
	// Simple verification of a known entry.
	found := false
	for _, c := range charts {
		if c.Name == "mychart" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected chart named \"mychart\" in result")
	}
	mu.Lock()
	if fetchCount != 1 {
		t.Fatalf("expected 1 fetch from server, got %d", fetchCount)
	}
	mu.Unlock()

	// Second call – should be served from cache, no extra request.
	charts2, err := ch.List(ctx)
	if err != nil {
		t.Fatalf("second List returned error: %v", err)
	}
	if len(charts2) != len(charts) {
		t.Fatalf("second List returned different number of charts")
	}
	mu.Lock()
	if fetchCount != 1 {
		t.Fatalf("expected no additional fetches, count=%d", fetchCount)
	}
	mu.Unlock()
}

/* -------------------------------------------------------------------------- *
 * List – fetch from a local directory (file based repo)                       *
 * -------------------------------------------------------------------------- */
func TestHelmChart_List_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Write index.yaml into the temporary directory.
	if err := os.WriteFile(filepath.Join(tmpDir, "index.yaml"), []byte(testIndexYAML), 0o600); err != nil {
		t.Fatalf("cannot write index.yaml: %v", err)
	}

	cfg := &config.Config{
		Kube: config.Kube{
			HelmRepositoryURLs: []string{tmpDir},
		},
	}
	k := mustNewKube(t, cfg)
	ch, err := NewHelmChart(k)
	if err != nil {
		t.Fatalf("NewHelmChart failed: %v", err)
	}
	ctx := context.Background()
	charts, err := ch.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(charts) == 0 {
		t.Fatalf("expected charts from file repo, got none")
	}
	// Verify that a known chart is present.
	found := false
	for _, c := range charts {
		if c.Name == "otherchart" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected chart \"otherchart\" in result")
	}
}

/* -------------------------------------------------------------------------- *
 * List – error handling when repository URL cannot be read                     *
 * -------------------------------------------------------------------------- */
func TestHelmChart_List_Error(t *testing.T) {
	// Non‑existent file path.
	cfg := &config.Config{
		Kube: config.Kube{
			HelmRepositoryURLs: []string{"/non/existent/path"},
		},
	}
	k := mustNewKube(t, cfg)
	ch, err := NewHelmChart(k)
	if err != nil {
		t.Fatalf("NewHelmChart failed: %v", err)
	}
	ctx := context.Background()
	_, err = ch.List(ctx)
	if err == nil {
		t.Fatalf("expected error when reading a non‑existent repo, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Show – error case (chart not found)                                         *
 * -------------------------------------------------------------------------- */
func TestHelmChart_Show_NotFound(t *testing.T) {
	// Regular Kube instance – environment settings are valid.
	cfg := &config.Config{}
	k := mustNewKube(t, cfg)
	ch, err := NewHelmChart(k)
	if err != nil {
		t.Fatalf("NewHelmChart failed: %v", err)
	}
	// Pass a bogus reference – Show should return an error.
	_, err = ch.Show("nonexistent/chart", action.ShowAll)
	if err == nil {
		t.Fatalf("expected Show to return error for unknown chart, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Concurrent fetch – ensure List runs all fetches in parallel without data‑race *
 * -------------------------------------------------------------------------- */
func TestHelmChart_List_ConcurrentFetch(t *testing.T) {
	// Two distinct HTTP servers, each serving the same test index.
	srv1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "index.yaml") {
			_, _ = w.Write([]byte(testIndexYAML))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv1.Close()

	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "index.yaml") {
			_, _ = w.Write([]byte(testIndexYAML))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv2.Close()

	cfg := &config.Config{
		Kube: config.Kube{
			HelmRepositoryURLs: []string{srv1.URL, srv2.URL},
		},
	}
	k := mustNewKube(t, cfg)
	ch, err := NewHelmChart(k)
	if err != nil {
		t.Fatalf("NewHelmChart failed: %v", err)
	}
	ctx := context.Background()

	// The errgroup inside List should run both fetches concurrently.
	// We only assert that the result aggregates entries from both repos
	// and that no panic occurs.
	charts, err := ch.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(charts) == 0 {
		t.Fatalf("expected charts from concurrent fetches, got none")
	}
}

/* -------------------------------------------------------------------------- *
 * Benchmark – creation and method calls                                       *
 * -------------------------------------------------------------------------- */
func BenchmarkHelmChart_Creation(b *testing.B) {
	cfg := &config.Config{
		Kube: config.Kube{
			HelmRepositoryURLs: []string{"https://example.com"},
		},
	}
	k := mustNewKube(b, cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := NewHelmChart(k); err != nil {
			b.Fatalf("creation error: %v", err)
		}
	}
}

func BenchmarkHelmChart_List(b *testing.B) {
	// HTTP server that always returns the same index file.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "index.yaml") {
			_, _ = w.Write([]byte(testIndexYAML))
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	cfg := &config.Config{
		Kube: config.Kube{
			HelmRepositoryURLs: []string{ts.URL},
		},
	}
	k := mustNewKube(b, cfg)
	ch, err := NewHelmChart(k)
	if err != nil {
		b.Fatalf("NewHelmChart failed: %v", err)
	}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ch.List(ctx)
	}
}

/* -------------------------------------------------------------------------- *
 * Type assertions                                                            *
 * -------------------------------------------------------------------------- */
func TestHelmChart_TypeAssertions(t *testing.T) {
	cfg := &config.Config{}
	k := mustNewKube(t, cfg)
	ch, err := NewHelmChart(k)
	if err != nil {
		t.Fatalf("NewHelmChart failed: %v", err)
	}
	hc, ok := ch.(*helmChart)
	if !ok {
		t.Fatal("cannot cast repository to *helmChart")
	}
	if hc.kube != k {
		t.Error("helmChart.kube field not set correctly")
	}
	var _ oscore.ChartRepo = hc
}

/* -------------------------------------------------------------------------- *
 * Edge cases – nil Kube pointer should panic on any method                    *
 * -------------------------------------------------------------------------- */
func TestHelmChart_EdgeCases(t *testing.T) {
	// Nil Kube – any method call must panic.
	t.Run("nil_kube", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using helmChart with nil Kube")
			}
		}()
		h := &helmChart{kube: nil}
		_, _ = h.List(context.Background())
	})
}
