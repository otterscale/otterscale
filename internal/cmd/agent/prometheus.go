package agent

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

// prometheusCAKey is where cert-manager writes the issuing CA in a
// Certificate's Secret.
const prometheusCAKey = "ca.crt"

// errPrometheusCA marks a CA that could not be loaded, so the error handler
// can say so without passing on the underlying error.
var errPrometheusCA = errors.New("cannot load the prometheus CA")

// caLoader returns the PEM bundle the Prometheus proxy's certificate must
// chain to.
type caLoader func(ctx context.Context) ([]byte, error)

// newPrometheusProxy forwards to target through rt. The caller's
// Authorization header is always dropped: the hub strips it already, and over
// https it would otherwise displace the agent's own token, which client-go's
// bearer round tripper leaves alone when a header is present. So are the
// Impersonate-* headers, which the hub strips too: the agent's ServiceAccount
// may impersonate, so nothing it sends may carry a caller's.
func newPrometheusProxy(target *url.URL, rt http.RoundTripper, errorLog *log.Logger) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.Header.Del("Authorization")
			// Header names arrive canonical, Impersonate-Extra-* included.
			for name := range r.Out.Header {
				if strings.HasPrefix(name, "Impersonate-") {
					r.Out.Header.Del(name)
				}
			}
		},
		Transport: rt,
		ModifyResponse: func(resp *http.Response) error {
			if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
				slog.Warn("prometheus proxy rejected the agent", "status", resp.StatusCode, "path", resp.Request.URL.Path)
			}
			return nil
		},
		ErrorHandler: prometheusProxyError,
		ErrorLog:     errorLog,
	}
}

// prometheusProxyError answers a request that never reached Prometheus with
// 502 and a reason the dashboard can show. The error itself, which can name
// Secrets and addresses, goes to the log only.
func prometheusProxyError(w http.ResponseWriter, r *http.Request, err error) {
	// A caller that went away cancels the request: nothing is wrong upstream.
	if r.Context().Err() != nil {
		slog.Debug("prometheus proxy request canceled by the caller", "path", r.URL.Path, "error", err)
	} else {
		slog.Warn("prometheus proxy error", "path", r.URL.Path, "error", err)
	}

	reason := "the agent cannot reach Prometheus"
	var verifyErr *tls.CertificateVerificationError
	switch {
	case errors.Is(err, errPrometheusCA):
		reason = "the agent cannot load the Prometheus proxy's CA"
	case errors.As(err, &verifyErr):
		reason = "the agent cannot verify the Prometheus proxy's certificate"
	}
	writePrometheusError(w, http.StatusBadGateway, reason)
}

// writePrometheusError answers in Prometheus' own error shape, as the hub
// does, so the dashboard shows the message rather than a generic error.
func writePrometheusError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		Status    string `json:"status"`
		ErrorType string `json:"errorType"`
		Error     string `json:"error"`
	}{"error", "unavailable", message})
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

const (
	// caRetryInterval is the least time between two CA loads once the last
	// one failed, or its CA failed verification. A missing Secret, or a URL
	// whose name the certificate does not carry, would otherwise cost an API
	// read and a handshake on every query.
	caRetryInterval = 15 * time.Second
	// caLoadTimeout bounds a load, which runs for every request waiting on it
	// rather than under any one request's context.
	caLoadTimeout = 10 * time.Second
)

// caTransport verifies the server against a CA read from the cluster. The CA
// is loaded on first use rather than at startup, because prometheus-stack is a
// module that may be installed after the agent, and loaded again when
// verification fails, so a re-issued CA is picked up without a restart.
type caTransport struct {
	loadCA caLoader
	now    func() time.Time // time.Now when nil

	loads singleflight.Group

	mu       sync.Mutex
	rt       *http.Transport
	bundle   []byte    // what rt trusts
	failure  error     // why rt cannot be used now; nil when it can
	retryAt  time.Time // while failure is set, no load before this
	loadedAt time.Time
}

func (t *caTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt, err := t.current(req.Context())
	if err != nil {
		return nil, err
	}

	resp, err := rt.RoundTrip(req)
	var verifyErr *tls.CertificateVerificationError
	if errors.As(err, &verifyErr) {
		t.fail(rt, err)
	}
	return resp, err
}

// current returns the transport to use, loading the CA when one is due. The
// load is shared by every request that needs it, and each of them can stop
// waiting on its own context.
func (t *caTransport) current(ctx context.Context) (*http.Transport, error) {
	if rt, due, err := t.cached(); !due {
		return rt, err
	}

	loaded := t.loads.DoChan("ca", func() (any, error) {
		loadCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), caLoadTimeout)
		defer cancel()
		return t.load(loadCtx)
	})
	select {
	case res := <-loaded:
		if res.Err != nil {
			return nil, res.Err
		}
		rt, ok := res.Val.(*http.Transport)
		if !ok {
			return nil, errors.New("the prometheus CA load returned no transport")
		}
		return rt, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// cached answers from what the last load left, unless a load is due: there
// is no transport yet, or the last failure is a caRetryInterval old.
func (t *caTransport) cached() (rt *http.Transport, due bool, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	switch {
	case t.rt != nil && t.failure == nil:
		return t.rt, false, nil
	case t.failure != nil && t.clock().Before(t.retryAt):
		return nil, false, t.failure
	}
	return nil, true, nil
}

// load reads the CA and keeps the current transport when the CA is the one
// it already trusts, so its connections are reused rather than redialled.
func (t *caTransport) load(ctx context.Context) (*http.Transport, error) {
	// A request that saw a load due may arrive after another load settled it.
	if rt, due, err := t.cached(); !due {
		return rt, err
	}

	bundle, err := t.loadCA(ctx)
	pool := x509.NewCertPool()
	if err == nil && !pool.AppendCertsFromPEM(bundle) {
		err = errors.New("the bundle holds no PEM certificate")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.loadedAt = t.clock()
	if err != nil {
		t.failure = fmt.Errorf("%w: %w", errPrometheusCA, err)
		t.retryAt = t.loadedAt.Add(caRetryInterval)
		return nil, t.failure
	}
	t.failure = nil
	if t.rt != nil && bytes.Equal(t.bundle, bundle) {
		return t.rt, nil
	}

	base, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.New("http.DefaultTransport is not an *http.Transport")
	}
	if t.rt != nil {
		t.rt.CloseIdleConnections()
	}
	t.rt = base.Clone()
	t.rt.TLSClientConfig = &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}
	t.bundle = bundle
	return t.rt, nil
}

// fail records that rt's CA failed verification, unless a load has replaced
// rt since. The next load is due a caRetryInterval after the last one, which
// a CA re-issued long after it is already past.
func (t *caTransport) fail(rt *http.Transport, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.rt != rt || t.failure != nil {
		return
	}
	t.failure = err
	t.retryAt = t.loadedAt.Add(caRetryInterval)
	rt.CloseIdleConnections()
}

func (t *caTransport) clock() time.Time {
	if t.now != nil {
		return t.now()
	}
	return time.Now()
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
