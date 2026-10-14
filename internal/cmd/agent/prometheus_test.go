package agent

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

func certPEM(der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

// unrelatedCA is a CA that did not sign httptest's certificate.
func unrelatedCA(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "unrelated"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return certPEM(der)
}

func tokenFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(path, []byte("agent-token"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func staticCA(pemBytes []byte) caLoader {
	return func(context.Context) ([]byte, error) { return pemBytes, nil }
}

func newTestProxy(t *testing.T, target string, cfg *rest.Config, loadCA caLoader) http.Handler {
	t.Helper()
	u, err := url.Parse(target)
	if err != nil {
		t.Fatal(err)
	}
	rt, err := newPrometheusTransport(u, cfg, loadCA)
	if err != nil {
		t.Fatal(err)
	}
	return newPrometheusProxy(u, rt, log.New(io.Discard, "", 0))
}

// query sends what the hub relays: a query carrying the caller's credential.
func query(h http.Handler) int {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/query?query=up", http.NoBody)
	req.Header.Set("Authorization", "Bearer caller-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

func TestPrometheusProxyPresentsAgentToken(t *testing.T) {
	var gotAuth, gotPath, gotQuery string
	srv := httptest.NewTLSServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		gotAuth, gotPath, gotQuery = r.Header.Get("Authorization"), r.URL.Path, r.URL.RawQuery
	}))
	defer srv.Close()

	cfg := &rest.Config{BearerTokenFile: tokenFile(t)}
	if code := query(newTestProxy(t, srv.URL, cfg, staticCA(certPEM(srv.Certificate().Raw)))); code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if gotAuth != "Bearer agent-token" {
		t.Errorf("upstream Authorization = %q, want the agent's token", gotAuth)
	}
	if gotPath != "/api/v1/query" || gotQuery != "query=up" {
		t.Errorf("upstream got %s?%s, want /api/v1/query?query=up", gotPath, gotQuery)
	}
}

func TestPrometheusProxyRefusesUnverifiedServer(t *testing.T) {
	reached := false
	srv := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { reached = true }))
	defer srv.Close()

	cfg := &rest.Config{BearerTokenFile: tokenFile(t)}
	if code := query(newTestProxy(t, srv.URL, cfg, staticCA(unrelatedCA(t)))); code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", code)
	}
	if reached {
		t.Fatal("the request, and the agent's token with it, reached a server the agent could not verify")
	}
}

// fakeClock is a caTransport's clock that moves only when told to.
type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func newFakeClock() *fakeClock { return &fakeClock{t: time.Unix(1_700_000_000, 0)} }

func (c *fakeClock) now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

// newCAProxy is newTestProxy with the caTransport given, for tests that need
// its clock or its state.
func newCAProxy(t *testing.T, target string, ct *caTransport) http.Handler {
	t.Helper()
	u, err := url.Parse(target)
	if err != nil {
		t.Fatal(err)
	}
	rt, err := transport.NewBearerAuthWithRefreshRoundTripper("", tokenFile(t), ct)
	if err != nil {
		t.Fatal(err)
	}
	return newPrometheusProxy(u, rt, log.New(io.Discard, "", 0))
}

// countLoads wraps answer and counts the calls.
func countLoads(answer caLoader) (caLoader, *atomic.Int32) {
	var n atomic.Int32
	return func(ctx context.Context) ([]byte, error) {
		n.Add(1)
		return answer(ctx)
	}, &n
}

// prometheus-stack may be installed after the agent, and its CA may be
// re-issued: neither should need an agent restart.
func TestPrometheusProxyRetriesCA(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()

	answers := []func() ([]byte, error){
		func() ([]byte, error) { return nil, errors.New("secret not found") },
		func() ([]byte, error) { return unrelatedCA(t), nil },
		func() ([]byte, error) { return certPEM(srv.Certificate().Raw), nil },
	}
	loads := 0
	loadCA := func(context.Context) ([]byte, error) {
		answer := answers[min(loads, len(answers)-1)]
		loads++
		return answer()
	}

	clock := newFakeClock()
	proxy := newCAProxy(t, srv.URL, &caTransport{loadCA: loadCA, now: clock.now})
	for i, want := range []int{http.StatusBadGateway, http.StatusBadGateway, http.StatusOK, http.StatusOK} {
		if code := query(proxy); code != want {
			t.Fatalf("request %d: status = %d, want %d", i+1, code, want)
		}
		clock.advance(caRetryInterval)
	}
	if loads != len(answers) {
		t.Errorf("CA loaded %d times, want %d: once per failure, then cached", loads, len(answers))
	}
}

// TestCATransportCoolsDownAfterFailedLoad: a CA that cannot be loaded is not
// asked for again on every query, only once per caRetryInterval.
func TestCATransportCoolsDownAfterFailedLoad(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()

	loadCA, loads := countLoads(func(context.Context) ([]byte, error) { return nil, errors.New("secret not found") })
	clock := newFakeClock()
	proxy := newCAProxy(t, srv.URL, &caTransport{loadCA: loadCA, now: clock.now})

	for i := range 3 {
		if code := query(proxy); code != http.StatusBadGateway {
			t.Fatalf("request %d: status = %d, want 502", i+1, code)
		}
	}
	if n := loads.Load(); n != 1 {
		t.Errorf("CA loaded %d times within caRetryInterval, want 1", n)
	}

	clock.advance(caRetryInterval)
	if code := query(proxy); code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", code)
	}
	if n := loads.Load(); n != 2 {
		t.Errorf("CA loaded %d times after caRetryInterval, want 2", n)
	}
}

// TestCATransportKeepsUnchangedCA: a certificate that fails verification
// with the right CA — here, a URL naming a host the certificate does not
// carry — costs one load and one handshake per caRetryInterval, and a reload
// that finds the same CA keeps the transport it has.
func TestCATransportKeepsUnchangedCA(t *testing.T) {
	var handshakes atomic.Int32
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	srv.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			handshakes.Add(1)
		}
	}
	srv.StartTLS()
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	misnamed := "https://localhost:" + u.Port() // the certificate names 127.0.0.1 and example.com

	loadCA, loads := countLoads(staticCA(certPEM(srv.Certificate().Raw)))
	clock := newFakeClock()
	ct := &caTransport{loadCA: loadCA, now: clock.now}
	proxy := newCAProxy(t, misnamed, ct)

	for i := range 3 {
		rec := serve(proxy, "/api/v1/query?query=up")
		assertPrometheusError(t, rec, http.StatusBadGateway, "cannot verify the Prometheus proxy's certificate")
		if t.Failed() {
			t.Fatalf("request %d", i+1)
		}
	}
	if n, h := loads.Load(), handshakes.Load(); n != 1 || h != 1 {
		t.Errorf("within caRetryInterval: %d loads and %d handshakes, want 1 and 1", n, h)
	}
	first := ct.rt

	clock.advance(caRetryInterval)
	if code := query(proxy); code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", code)
	}
	if n, h := loads.Load(), handshakes.Load(); n != 2 || h != 2 {
		t.Errorf("after caRetryInterval: %d loads and %d handshakes, want 2 and 2", n, h)
	}
	if ct.rt != first {
		t.Error("the transport was rebuilt for the same CA")
	}
}

// TestCATransportWaitersCanGiveUp: requests waiting on a CA load share it,
// and each can stop waiting when its own context ends.
func TestCATransportWaitersCanGiveUp(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()

	started, release := make(chan struct{}), make(chan struct{})
	loadCA, loads := countLoads(func(context.Context) ([]byte, error) {
		close(started)
		<-release
		return certPEM(srv.Certificate().Raw), nil
	})
	proxy := newCAProxy(t, srv.URL, &caTransport{loadCA: loadCA})

	first := make(chan int, 1)
	go func() { first <- query(proxy) }()
	<-started

	ctx, cancel := context.WithCancel(context.Background())
	waiter := make(chan int, 1)
	go func() {
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/query?query=up", http.NoBody)
		rec := httptest.NewRecorder()
		proxy.ServeHTTP(rec, req)
		waiter <- rec.Code
	}()
	cancel()
	select {
	case <-waiter:
	case <-time.After(5 * time.Second):
		t.Fatal("a waiting request did not stop when its context ended")
	}

	close(release)
	if code := <-first; code != http.StatusOK {
		t.Errorf("first request: status = %d, want 200", code)
	}
	if n := loads.Load(); n != 1 {
		t.Errorf("CA loaded %d times, want 1 shared load", n)
	}
}

func TestPrometheusProxyPassesRejectionThrough(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	cfg := &rest.Config{BearerTokenFile: tokenFile(t)}
	if code := query(newTestProxy(t, srv.URL, cfg, staticCA(certPEM(srv.Certificate().Raw)))); code != http.StatusForbidden {
		t.Fatalf("status = %d, want the upstream 403", code)
	}
}

// assertPrometheusError requires Prometheus' JSON error shape, whose message
// the dashboard shows, and that message to contain want.
func assertPrometheusError(t *testing.T, rec *httptest.ResponseRecorder, status int, want string) {
	t.Helper()
	if rec.Code != status {
		t.Errorf("status = %d, want %d", rec.Code, status)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body struct {
		Status    string `json:"status"`
		ErrorType string `json:"errorType"`
		Error     string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil ||
		body.Status != "error" || body.ErrorType != "unavailable" || !strings.Contains(body.Error, want) {
		t.Errorf("body = %q, want Prometheus' error shape with %q in it", rec.Body, want)
	}
}

func serve(h http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// TestPrometheusProxyExplainsFailures: a request the agent cannot deliver is
// 502 with a reason the dashboard can show, and nothing from the underlying
// error — Secret names, addresses — in it.
func TestPrometheusProxyExplainsFailures(t *testing.T) {
	tlsSrv := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer tlsSrv.Close()
	closed := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	closed.Close()

	cfg := &rest.Config{BearerTokenFile: tokenFile(t)}
	cases := []struct {
		desc   string
		target string
		loadCA caLoader
		want   string
	}{
		{
			"CA not loaded", tlsSrv.URL,
			func(context.Context) ([]byte, error) { return nil, errors.New("secrets \"prom-tls\" not found") },
			"cannot load the Prometheus proxy's CA",
		},
		{"CA holds no certificate", tlsSrv.URL, staticCA([]byte("not PEM")), "cannot load the Prometheus proxy's CA"},
		{"certificate not verified", tlsSrv.URL, staticCA(unrelatedCA(t)), "cannot verify the Prometheus proxy's certificate"},
		{"nothing listening", closed.URL, nil, "cannot reach Prometheus"},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			rec := serve(newTestProxy(t, c.target, cfg, c.loadCA), "/api/v1/query?query=up")
			assertPrometheusError(t, rec, http.StatusBadGateway, c.want)
			for _, secret := range []string{"prom-tls", "127.0.0.1"} {
				if strings.Contains(rec.Body.String(), secret) {
					t.Errorf("body %q carries %q", rec.Body, secret)
				}
			}
		})
	}
}

// TestMountDisabledPrometheus: with the metrics proxy off, its prefix answers
// 503 with the reason instead of falling through to kube-apiserver's proxy.
func TestMountDisabledPrometheus(t *testing.T) {
	kubeCalled := false
	kube := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		kubeCalled = true
		w.WriteHeader(http.StatusNotFound)
	}))
	defer kube.Close()

	cases := []struct {
		desc string
		url  string
		want string
	}{
		{"no ServiceAccount token", "https://prometheus.example:9091", "no ServiceAccount token"},
		{"no Prometheus URL", "", "no Prometheus URL"},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			h := &Handler{
				cfg:                     &rest.Config{Host: kube.URL},
				proxyPrometheusURL:      c.url,
				proxyPrometheusCASecret: "monitoring/prom-tls",
			}
			mux := http.NewServeMux()
			if err := h.Mount(mux); err != nil {
				t.Fatalf("Mount() = %v, want nil: a disabled metrics proxy must not stop the agent", err)
			}

			assertPrometheusError(t, serve(mux, proxyPathPrefix+"api/v1/query?query=up"), http.StatusServiceUnavailable, c.want)
			if kubeCalled {
				t.Error("the request fell through to kube-apiserver's proxy")
			}
		})
	}
}

func TestPrometheusProxyHTTPWithholdsToken(t *testing.T) {
	gotAuth, reached := "unset", false
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		gotAuth, reached = r.Header.Get("Authorization"), true
	}))
	defer srv.Close()

	cfg := &rest.Config{BearerTokenFile: tokenFile(t)}
	if code := query(newTestProxy(t, srv.URL, cfg, nil)); code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if !reached || gotAuth != "" {
		t.Errorf("upstream Authorization = %q, want none over plain http", gotAuth)
	}
}

func TestPrometheusTransportNeedsToken(t *testing.T) {
	u, err := url.Parse("https://prometheus.example:9091")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := newPrometheusTransport(u, &rest.Config{}, staticCA(nil)); err == nil {
		t.Fatal("want an error for an https target without a ServiceAccount token")
	}
}

func TestPrometheusTransportRejectsOtherSchemes(t *testing.T) {
	u, err := url.Parse("ftp://prometheus.example:9091")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := newPrometheusTransport(u, &rest.Config{BearerToken: "t"}, nil); err == nil {
		t.Fatal("want an error for a non-http(s) target")
	}
}

func TestSecretCALoader(t *testing.T) {
	client := fake.NewClientset(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Name: "prom-tls"},
			Data:       map[string][]byte{"ca.crt": []byte("CA"), "tls.key": []byte("KEY")},
		},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Name: "no-ca"}},
	)

	load, err := secretCALoader(client, "monitoring/prom-tls")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := load(context.Background()); err != nil || string(got) != "CA" {
		t.Errorf("load() = %q, %v; want ca.crt", got, err)
	}

	for _, ref := range []string{"monitoring/missing", "monitoring/no-ca"} {
		load, err := secretCALoader(client, ref)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := load(context.Background()); err == nil {
			t.Errorf("%s: want an error", ref)
		}
	}

	for _, ref := range []string{"prom-tls", "/prom-tls", "monitoring/", "a/b/c", ""} {
		if _, err := secretCALoader(client, ref); err == nil {
			t.Errorf("%q: want an error for a ref that is not namespace/name", ref)
		}
	}
}
