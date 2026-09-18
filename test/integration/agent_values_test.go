package integration

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"connectrpc.com/authn"
	"connectrpc.com/connect"

	linkv1 "github.com/otterscale/otterscale/api/link/v1"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/handler"
	"github.com/otterscale/otterscale/internal/providers/cache"
	"github.com/otterscale/otterscale/internal/providers/values"
	otterhttp "github.com/otterscale/otterscale/internal/transport/http"
)

// Mirrors the route the server mounts, repeated rather than exported so that
// changing it there without changing it here fails this test.
const agentValuesPath = "/link/values/"

// Issuing values contacts Harbor; fetching them through the URL must not,
// which this counts.
type stubHarbor struct{ calls int }

func (h *stubHarbor) EnsureRobotAccount(
	_ context.Context, cluster, secret string,
) (core.HarborRobotCredentials, error) {
	h.calls++
	return core.HarborRobotCredentials{Name: "robot$" + cluster, Secret: secret}, nil
}

// The test is about the public route, so the OIDC middleware is out of scope.
func adminAuth() *authn.Middleware {
	return authn.NewMiddleware(func(_ context.Context, r *http.Request) (any, error) {
		if r.Header.Get("Authorization") == "" {
			return nil, authn.Errorf("missing bearer token")
		}
		return core.UserInfo{
			Subject: "2ae606b7-ff7c-4c47-be76-8371963ee3c5",
			Groups:  []string{"system:authenticated", "oidc:admin"},
		}, nil
	})
}

// The end-to-end form of the property the whole design turns on: an operator
// runs the RPC and then fetches the URL it returned, so the two have to be the
// same file. The implementation this replaces revoked the first credential on
// every fetch.
func TestAgentValuesRPCAndSignedURLServeIdenticalBytes(t *testing.T) {
	join, err := core.NewJoinAuthority(integrationSecret)
	if err != nil {
		t.Fatalf("NewJoinAuthority: %v", err)
	}
	// The listener comes first, so the URL the RPC returns can be fetched as
	// given, exactly as an operator would.
	var lc net.ListenConfig
	ln, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	baseURL := "http://" + ln.Addr().String()

	harbor := &stubHarbor{}
	useCase := core.NewAgentValuesUseCase(
		&core.AgentValuesConfig{
			ExternalURL:     baseURL,
			TunnelServerURL: "https://192.0.2.1:30300",
			HarborURL:       "https://harbor.example.com:8443",
			TrustedCASecret: core.TrustedCASecretName,
			TrustedCAKey:    core.DefaultTrustedCAKey,
		},
		core.Version("test"), join, cache.NewAgentValuesStore(), values.NewRenderer(), harbor,
	)

	linkService := handler.NewLinkService(core.NewLinkUseCase(newTestTunnel(t), "test", join), useCase)
	valuesHandler := handler.NewAgentValuesHandler(useCase)

	srv, err := otterhttp.NewServer(t.Context(),
		otterhttp.WithListener(ln),
		otterhttp.WithAuthMiddleware(adminAuth()),
		otterhttp.WithAllowedOrigins([]string{"https://example.com"}),
		otterhttp.WithPublicPathPrefixes([]string{agentValuesPath}),
		otterhttp.WithMount(func(mux *http.ServeMux) error {
			mux.Handle(linkv1.NewLinkServiceHandler(linkService))
			mux.HandleFunc("GET "+agentValuesPath+"{id}", func(w http.ResponseWriter, r *http.Request) {
				rendered, err := valuesHandler.Render(r.Context(), r.PathValue("id"))
				if err != nil {
					http.Error(w, "invalid or expired token", http.StatusUnauthorized)
					return
				}
				w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
				w.Header().Set("Cache-Control", "no-store")
				_, _ = w.Write([]byte(rendered)) //nolint:gosec // G705: YAML served as text/yaml, not markup
			})
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}

	httpSrv := httptest.NewUnstartedServer(srv.Handler())
	httpSrv.Listener.Close()
	httpSrv.Listener = ln
	httpSrv.Start()
	t.Cleanup(httpSrv.Close)

	client := linkv1.NewLinkServiceClient(httpSrv.Client(), baseURL,
		connect.WithInterceptors(bearerTokenInterceptor("test-token")))

	info := &linkv1.AgentClusterInfo{}
	info.SetExternalAddress("192.0.2.1")
	info.SetNodePortRange("30000-32767")

	req := &linkv1.IssueAgentValuesRequest{}
	req.SetCluster("prod")
	req.SetClusterInfo(info)

	resp, err := client.IssueAgentValues(t.Context(), req)
	if err != nil {
		t.Fatalf("IssueAgentValues: %v", err)
	}
	if resp.GetValues() == "" {
		t.Fatal("the RPC returned no values")
	}
	if harbor.calls != 1 {
		t.Fatalf("harbor calls after the RPC = %d, want 1", harbor.calls)
	}

	// Fetched exactly as returned, with no rewriting.
	fetched := fetchValues(t, resp.GetUrl())

	if fetched != resp.GetValues() {
		t.Errorf("the URL served different bytes than the RPC:\n%s\n---\n%s", fetched, resp.GetValues())
	}
	if harbor.calls != 1 {
		t.Errorf("harbor calls after the fetch = %d, want 1: the public route must not write to Harbor", harbor.calls)
	}

	// Fetching again must be free of consequence: Helm retries, Flux
	// reconciles, operators re-run installs.
	if again := fetchValues(t, resp.GetUrl()); again != fetched {
		t.Error("two fetches of the same URL returned different bytes")
	}
	if harbor.calls != 1 {
		t.Errorf("harbor calls after two fetches = %d, want 1", harbor.calls)
	}
}

// The generated client offers no other way to set a header.
func bearerTokenInterceptor(token string) connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			req.Header().Set("Authorization", "Bearer "+token)
			return next(ctx, req)
		}
	})
}

func fetchValues(t *testing.T, url string) string {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	// Deliberately no Authorization header: the id in the path is the whole
	// credential.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d (body %q)", resp.StatusCode, http.StatusOK, body)
	}
	if got, want := resp.Header.Get("Cache-Control"), "no-store"; got != want {
		t.Errorf("Cache-Control = %q, want %q", got, want)
	}
	if got := resp.Header.Get("Content-Type"); !strings.HasPrefix(got, "text/yaml") {
		t.Errorf("Content-Type = %q, want text/yaml", got)
	}

	return string(body)
}
