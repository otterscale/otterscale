package http

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	var lc net.ListenConfig
	ln, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	srv, err := NewServer(
		t.Context(),
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
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/public", http.NoBody)
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("private path without token is blocked", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/private", http.NoBody)
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)

		if rec.Code == http.StatusOK {
			t.Fatalf("expected non-200 status for private path without token, got %d", rec.Code)
		}
	})

	t.Run("private path with token is allowed", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/private", http.NoBody)
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})
}

// newTestServer builds a Server on an ephemeral loopback listener.
func newTestServer(t *testing.T, opts ...ServerOption) *Server {
	t.Helper()

	var lc net.ListenConfig
	ln, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	srv, err := NewServer(t.Context(), append([]ServerOption{WithListener(ln)}, opts...)...)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	return srv
}

// deadlineRecorder records the deadlines http.ResponseController sets
// on it. ResponseController looks these methods up on the
// ResponseWriter, so a plain httptest.ResponseRecorder cannot observe
// them.
type deadlineRecorder struct {
	http.ResponseWriter

	readDeadline  *time.Time
	writeDeadline *time.Time
}

func (d *deadlineRecorder) SetReadDeadline(deadline time.Time) error {
	d.readDeadline = &deadline
	return nil
}

func (d *deadlineRecorder) SetWriteDeadline(deadline time.Time) error {
	d.writeDeadline = &deadline
	return nil
}

func serveRecordingDeadlines(t *testing.T, srv *Server, path string) *deadlineRecorder {
	t.Helper()

	rec := &deadlineRecorder{ResponseWriter: httptest.NewRecorder()}
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, http.NoBody)
	srv.Handler().ServeHTTP(rec, req)
	return rec
}

func okMount(paths ...string) MountFunc {
	return func(mux *http.ServeMux) error {
		for _, p := range paths {
			mux.HandleFunc(p, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
		}
		return nil
	}
}

// TestServer_LongRunningPathClearsDeadlines guards the streaming RPCs
// against the write timeout, which over HTTP/1.1 bounds the whole
// response and would otherwise end every watch or exec session after
// requestTimeout.
func TestServer_LongRunningPathClearsDeadlines(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t,
		WithLongRunningPaths([]string{"/stream"}),
		WithMount(okMount("/stream", "/unary")),
	)

	rec := serveRecordingDeadlines(t, srv, "/stream")

	if rec.readDeadline == nil || !rec.readDeadline.IsZero() {
		t.Errorf("read deadline = %v, want the zero time (no deadline)", rec.readDeadline)
	}
	if rec.writeDeadline == nil || !rec.writeDeadline.IsZero() {
		t.Errorf("write deadline = %v, want the zero time (no deadline)", rec.writeDeadline)
	}
}

func TestServer_RegularPathKeepsDeadlines(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t,
		WithLongRunningPaths([]string{"/stream"}),
		WithMount(okMount("/stream", "/unary")),
	)

	rec := serveRecordingDeadlines(t, srv, "/unary")

	if rec.readDeadline != nil || rec.writeDeadline != nil {
		t.Errorf("deadlines were modified for a unary path: read=%v write=%v", rec.readDeadline, rec.writeDeadline)
	}
}

func TestServer_Timeouts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		opts      []ServerOption
		wantRead  time.Duration
		wantWrite time.Duration
	}{
		{
			name:      "defaults bound a request",
			wantRead:  requestTimeout,
			wantWrite: requestTimeout,
		},
		{
			name:      "agent disables them",
			opts:      []ServerOption{WithoutRequestTimeouts()},
			wantRead:  0,
			wantWrite: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := newTestServer(t, append(tt.opts, WithMount(okMount("/unary")))...)

			if got := srv.inner.ReadTimeout; got != tt.wantRead {
				t.Errorf("ReadTimeout = %v, want %v", got, tt.wantRead)
			}
			if got := srv.inner.WriteTimeout; got != tt.wantWrite {
				t.Errorf("WriteTimeout = %v, want %v", got, tt.wantWrite)
			}
			// These stay in place either way: headers are always
			// bounded, and idle keep-alive connections must not
			// outlive the request timeouts being disabled.
			if got := srv.inner.ReadHeaderTimeout; got != readHeaderTimeout {
				t.Errorf("ReadHeaderTimeout = %v, want %v", got, readHeaderTimeout)
			}
			if got := srv.inner.IdleTimeout; got != idleTimeout {
				t.Errorf("IdleTimeout = %v, want %v", got, idleTimeout)
			}
		})
	}
}
