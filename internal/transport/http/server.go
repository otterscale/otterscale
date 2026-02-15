package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/authn"
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"

	"github.com/otterscale/otterscale/internal/core"
)

// MountFunc registers handlers onto the provided ServeMux.
// Accepting *http.ServeMux allows the caller to register multiple services.
type MountFunc func(mux *http.ServeMux) error

// ServerOption configures a Server.
type ServerOption func(*Server)

// Server is an HTTP/H2C server with optional CORS and authentication
// middleware. It implements transport.Listener.
type Server struct {
	inner              *http.Server
	address            string
	listener           net.Listener
	mount              MountFunc
	authMiddleware     *authn.Middleware
	publicPaths        map[string]struct{}
	publicPathPrefixes []string
	allowedOrigins     []string
	log                *slog.Logger
}

// WithAddress configures the listen address (e.g. ":8299").
func WithAddress(address string) ServerOption {
	return func(s *Server) { s.address = address }
}

// WithListener provides an external net.Listener for the server to
// use. When set, Start will serve on this listener instead of
// creating a new TCP listener from the configured address.
func WithListener(ln net.Listener) ServerOption {
	return func(s *Server) { s.listener = ln }
}

// WithMount configures the function that registers route handlers.
func WithMount(mount MountFunc) ServerOption {
	return func(s *Server) { s.mount = mount }
}

// WithAuthMiddleware configures the authentication middleware.
func WithAuthMiddleware(m *authn.Middleware) ServerOption {
	return func(s *Server) { s.authMiddleware = m }
}

// WithPublicPaths configures paths that bypass authentication.
// Paths are normalised to always include a leading "/".
func WithPublicPaths(paths []string) ServerOption {
	return func(s *Server) {
		if len(paths) == 0 {
			return
		}
		if s.publicPaths == nil {
			s.publicPaths = make(map[string]struct{}, len(paths))
		}
		for _, p := range paths {
			if p == "" {
				continue
			}
			if p[0] != '/' {
				p = "/" + p
			}
			s.publicPaths[p] = struct{}{}
		}
	}
}

// WithPublicPathPrefixes configures path prefixes that bypass
// authentication. Any request whose path starts with one of these
// prefixes is served without OIDC token verification. Prefixes are
// normalised to always include a leading "/".
func WithPublicPathPrefixes(prefixes []string) ServerOption {
	return func(s *Server) {
		for _, p := range prefixes {
			if p == "" {
				continue
			}
			if p[0] != '/' {
				p = "/" + p
			}
			s.publicPathPrefixes = append(s.publicPathPrefixes, p)
		}
	}
}

// WithAllowedOrigins configures the allowed origins for CORS.
func WithAllowedOrigins(origins []string) ServerOption {
	return func(s *Server) { s.allowedOrigins = origins }
}

// WithHTTPLogger configures a structured logger. Defaults to
// slog.Default with a "component" attribute.
func WithHTTPLogger(log *slog.Logger) ServerOption {
	return func(s *Server) { s.log = log }
}

// NewServer creates a new HTTP server with the given options.
func NewServer(ctx context.Context, opts ...ServerOption) (*Server, error) {
	s := &Server{
		address: ":8299",
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.log == nil {
		s.log = slog.Default().With("component", "http-server")
	}
	// When authentication is enabled (server mode), require explicit
	// CORS origins to avoid accidentally exposing the API to all
	// origins in production.
	if s.authMiddleware != nil && len(s.allowedOrigins) == 0 {
		return nil, fmt.Errorf("http server: allowed origins must be configured when authentication is enabled; " +
			"set --allowed-origins or OTTERSCALE_SERVER_ALLOWED_ORIGINS")
	}
	if s.listener == nil {
		var lc net.ListenConfig
		ln, err := lc.Listen(ctx, "tcp", s.address)
		if err != nil {
			return nil, fmt.Errorf("http listen %q: %w", s.address, err)
		}
		s.listener = ln
	}

	handler, err := s.buildHandler()
	if err != nil {
		return nil, err
	}

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	s.inner = &http.Server{
		Addr:              s.address,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8 KiB
		Protocols:         protocols,
	}

	return s, nil
}

// Handler returns the server's top-level HTTP handler. This is useful
// for testing the middleware chain without starting a real listener.
func (s *Server) Handler() http.Handler {
	return s.inner.Handler
}

// Start begins accepting connections and blocks until the server is
// shut down or an unrecoverable error occurs.
func (s *Server) Start(ctx context.Context) error {
	s.inner.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	s.log.Info("starting",
		"address", s.listener.Addr().String(),
		"auth", s.authMiddleware != nil,
		"public_paths", len(s.publicPaths),
		"allowed_origins", s.allowedOrigins,
	)

	if err := s.inner.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http serve: %w", err)
	}

	return nil
}

// Stop gracefully drains connections. If the graceful shutdown
// exceeds the context deadline it forces an immediate close.
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("shutting down")
	if err := s.inner.Shutdown(ctx); err != nil {
		s.log.Error("graceful shutdown failed, forcing close", "error", err)
		return s.inner.Close()
	}
	return nil
}

// ---------------------------------------------------------------------------
// Middleware chain
// ---------------------------------------------------------------------------

// buildHandler assembles the middleware stack.
// Order: H2C -> CORS -> Auth -> Mux
func (s *Server) buildHandler() (http.Handler, error) {
	mux := http.NewServeMux()
	if s.mount != nil {
		if err := s.mount(mux); err != nil {
			return nil, fmt.Errorf("mount routes: %w", err)
		}
	}

	var handler http.Handler = mux

	// Authentication
	if s.authMiddleware != nil {
		handler = s.wrapAuth(mux, handler)
	}

	// CORS
	handler = s.wrapCORS(handler)

	return handler, nil
}

// wrapAuth applies the authn middleware, skipping public paths.
// Public paths are checked by exact match first, then by prefix.
// After authn sets the transport-level auth info, bridgeUserInfo
// copies it into the domain-level core.UserInfo context key so that
// infrastructure adapters can access the user identity without
// depending on the connectrpc/authn package.
func (s *Server) wrapAuth(mux *http.ServeMux, next http.Handler) http.Handler {
	bridged := bridgeUserInfo(next)
	protected := s.authMiddleware.Wrap(bridged)
	if len(s.publicPaths) == 0 && len(s.publicPathPrefixes) == 0 {
		return protected
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.isPublicPath(r.URL.Path) {
			mux.ServeHTTP(w, r)
			return
		}
		protected.ServeHTTP(w, r)
	})
}

// bridgeUserInfo extracts the authn-stored UserInfo and stores it via
// the domain-level core.WithUserInfo context accessor. This decouples
// infrastructure adapters from the transport-specific authn package.
func bridgeUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if info, ok := authn.GetInfo(r.Context()).(core.UserInfo); ok {
			r = r.WithContext(core.WithUserInfo(r.Context(), info))
		}
		next.ServeHTTP(w, r)
	})
}

// isPublicPath returns true if the given path matches an exact public
// path or starts with a registered public path prefix.
func (s *Server) isPublicPath(path string) bool {
	if _, ok := s.publicPaths[path]; ok {
		return true
	}
	for _, prefix := range s.publicPathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// wrapCORS applies CORS headers. When no origins are configured
// (agent mode) it allows all origins. This is safe because the agent
// serves exclusively on an in-memory pipe listener behind the chisel
// tunnel — traffic never reaches the agent directly from a browser.
// All requests are forwarded through the server's mTLS-authenticated
// tunnel, so browser-origin restrictions are enforced at the server
// layer instead. In server mode the startup validation in NewServer
// ensures allowedOrigins is non-empty.
func (s *Server) wrapCORS(next http.Handler) http.Handler {
	if len(s.allowedOrigins) == 0 {
		return cors.AllowAll().Handler(next)
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   s.allowedOrigins,
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   connectcors.AllowedHeaders(),
		ExposedHeaders:   connectcors.ExposedHeaders(),
		AllowCredentials: true,
		MaxAge:           7200,
	})
	return c.Handler(next)
}
