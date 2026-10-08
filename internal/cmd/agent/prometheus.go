package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

// prometheusCAKey is where cert-manager writes the issuing CA in a
// Certificate's Secret.
const prometheusCAKey = "ca.crt"

// caLoader returns the PEM bundle the Prometheus proxy's certificate must
// chain to.
type caLoader func(ctx context.Context) ([]byte, error)

// newPrometheusProxy forwards to target through rt. The caller's
// Authorization header is always dropped: the hub strips it already, and over
// https it would otherwise displace the agent's own token, which client-go's
// bearer round tripper leaves alone when a header is present.
func newPrometheusProxy(target *url.URL, rt http.RoundTripper, errorLog *log.Logger) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.Header.Del("Authorization")
		},
		Transport: rt,
		ErrorLog:  errorLog,
	}
}

// newPrometheusTransport picks how the agent talks to Prometheus by the
// target's scheme.
//
// https is kube-rbac-proxy: every request carries the agent's ServiceAccount
// token, re-read from cfg's token file so a rotated token is picked up, and
// the server must chain to the CA loadCA returns. The token can impersonate
// any user, so it is never sent to a server the agent has not verified.
//
// http is a Prometheus without the proxy, as an older prometheus-stack
// deploys it: nothing is added, so the token never crosses the network in
// clear.
func newPrometheusTransport(target *url.URL, cfg *rest.Config, loadCA caLoader) (http.RoundTripper, error) {
	switch target.Scheme {
	case "http":
		return http.DefaultTransport, nil
	case "https":
		if cfg.BearerToken == "" && cfg.BearerTokenFile == "" {
			return nil, errors.New("an https prometheus URL is kube-rbac-proxy, which needs the agent's ServiceAccount token, and this kube config carries none")
		}
		return transport.NewBearerAuthWithRefreshRoundTripper(cfg.BearerToken, cfg.BearerTokenFile, &caTransport{loadCA: loadCA})
	default:
		return nil, fmt.Errorf("prometheus URL scheme must be http or https, got %q", target.Scheme)
	}
}

// caTransport verifies the server against a CA read from the cluster. The CA
// is loaded on first use rather than at startup, because prometheus-stack is a
// module that may be installed after the agent, and it is dropped whenever
// verification fails, so a re-issued CA is picked up without a restart.
type caTransport struct {
	loadCA caLoader

	mu sync.Mutex
	rt *http.Transport
}

func (t *caTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt, err := t.current(req.Context())
	if err != nil {
		return nil, err
	}

	resp, err := rt.RoundTrip(req)
	var verifyErr *tls.CertificateVerificationError
	if errors.As(err, &verifyErr) {
		t.drop(rt)
	}
	return resp, err
}

func (t *caTransport) current(ctx context.Context) (*http.Transport, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.rt != nil {
		return t.rt, nil
	}

	bundle, err := t.loadCA(ctx)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(bundle) {
		return nil, errors.New("the prometheus CA bundle holds no PEM certificate")
	}

	base, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.New("http.DefaultTransport is not an *http.Transport")
	}
	rt := base.Clone()
	rt.TLSClientConfig = &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}
	t.rt = rt
	return rt, nil
}

// drop forgets rt, unless a concurrent request already replaced it.
func (t *caTransport) drop(rt *http.Transport) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.rt == rt {
		t.rt = nil
		rt.CloseIdleConnections()
	}
}

// secretCALoader reads ca.crt from the Secret ref names as namespace/name.
func secretCALoader(client kubernetes.Interface, ref string) (caLoader, error) {
	namespace, name, ok := strings.Cut(ref, "/")
	if !ok || namespace == "" || name == "" || strings.Contains(name, "/") {
		return nil, fmt.Errorf("prometheus CA secret must be namespace/name, got %q", ref)
	}

	return func(ctx context.Context) ([]byte, error) {
		secret, err := client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to read the prometheus CA from secret %s: %w", ref, err)
		}
		bundle := secret.Data[prometheusCAKey]
		if len(bundle) == 0 {
			return nil, fmt.Errorf("secret %s has no %s", ref, prometheusCAKey)
		}
		return bundle, nil
	}, nil
}
