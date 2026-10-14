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
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
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

	cfg := &rest.Config{BearerTokenFile: tokenFile(t)}
	proxy := newTestProxy(t, srv.URL, cfg, loadCA)
	for i, want := range []int{http.StatusBadGateway, http.StatusBadGateway, http.StatusOK, http.StatusOK} {
		if code := query(proxy); code != want {
			t.Fatalf("request %d: status = %d, want %d", i+1, code, want)
		}
	}
	if loads != len(answers) {
		t.Errorf("CA loaded %d times, want %d: once per failure, then cached", loads, len(answers))
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
