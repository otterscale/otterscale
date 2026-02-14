package http

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/authn"
)

func TestNewServer_PublicPathsBypassAuth(t *testing.T) {
	t.Parallel()

	authMiddleware := authn.NewMiddleware(func(_ context.Context, r *http.Request) (any, error) {
		if r.Header.Get("Authorization") == "" {
			return nil, authn.Errorf("missing bearer token")
		}
		return struct{}{}, nil
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	srv, err := NewServer(
		WithListener(ln),
		WithAuthMiddleware(authMiddleware),
		WithAllowedOrigins([]string{"https://example.com"}),
		WithPublicPaths([]string{"/public"}),
		WithMount(func(mux *http.ServeMux) error {
			mux.HandleFunc("/public", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			mux.HandleFunc("/private", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	t.Run("public path without token is allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("private path without token is blocked", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/private", nil)
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)

		if rec.Code == http.StatusOK {
			t.Fatalf("expected non-200 status for private path without token, got %d", rec.Code)
		}
	})

	t.Run("private path with token is allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/private", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})
}
