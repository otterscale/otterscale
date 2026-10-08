package agent

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
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
