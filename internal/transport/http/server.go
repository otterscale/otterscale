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

// Request timeouts. ReadHeaderTimeout always applies; the read and
// write timeouts bound a single request/response exchange and are
// therefore lifted per-request for long-running paths (see
// WithLongRunningPaths) and disabled entirely by WithoutRequestTimeouts.
// IdleTimeout is set explicitly rather than inherited from ReadTimeout
// so that disabling the latter does not leave keep-alive connections
// without any bound at all.
const (
	readHeaderTimeout = 5 * time.Second
	requestTimeout    = 5 * time.Minute
	idleTimeout       = 5 * time.Minute
)

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
	longRunningPaths   map[string]struct{}
	noRequestTimeouts  bool
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

// WithLongRunningPaths marks request paths whose response is a
// long-lived stream (server-streaming RPCs such as watch, log follow,
// exec, port-forward and VNC). Requests to those paths have their
// connection deadlines cleared, because the request timeouts bound a
// whole exchange: over HTTP/1.1 they would otherwise cut every such
// stream off after requestTimeout.
func WithLongRunningPaths(paths []string) ServerOption {
	return func(s *Server) {
		if s.longRunningPaths == nil {
			s.longRunningPaths = make(map[string]struct{}, len(paths))
		}
		for _, p := range paths {
			if p == "" {
				continue
			}
			if p[0] != '/' {
				p = "/" + p
			}
			s.longRunningPaths[p] = struct{}{}
		}
	}
}

// WithoutRequestTimeouts disables the read and write timeouts. It is
// meant for the agent, which proxies arbitrary kube-apiserver traffic
// (exec, attach, port-forward, log follow, watch) whose duration is
// unbounded and whose upgraded connections keep the deadlines that were
// set before the connection was hijacked. The agent serves exclusively
// on an in-memory pipe behind the tunnel, so there is no untrusted
// network peer these timeouts would protect it from.
func WithoutRequestTimeouts() ServerOption {
	return func(s *Server) { s.noRequestTimeouts = true }
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

	read, write := requestTimeout, requestTimeout
	if s.noRequestTimeouts {
		read, write = 0, 0
	}

	s.inner = &http.Server{
		Addr:              s.address,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       read,
		WriteTimeout:      write,
		IdleTimeout:       idleTimeout,
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
// Order: H2C -> CORS -> Auth -> Deadlines -> Mux
func (s *Server) buildHandler() (http.Handler, error) {
	mux := http.NewServeMux()
	if s.mount != nil {
		if err := s.mount(mux); err != nil {
			return nil, fmt.Errorf("mount routes: %w", err)
		}
	}

	// Deadlines are lifted closest to the mux so that the decision is
	// made once per request, after routing information is available and
	// before any handler starts writing.
	handler := s.wrapDeadlines(mux)

	// Authentication
	if s.authMiddleware != nil {
		handler = s.wrapAuth(handler)
	}

	// CORS
	handler = s.wrapCORS(handler)

	return handler, nil
}

// wrapDeadlines clears the connection deadlines for long-running paths
// so that a streaming response is not cut off mid-flight.
//
// Over HTTP/1.1 the write timeout is a hard deadline on the entire
// response, so a watch or log-follow stream would end after
// requestTimeout. Over HTTP/2 the timeouts are progress-based and this
// is a no-op; SetWriteDeadline reporting ErrNotSupported is therefore
// not an error worth surfacing.
func (s *Server) wrapDeadlines(next http.Handler) http.Handler {
	if len(s.longRunningPaths) == 0 {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := s.longRunningPaths[r.URL.Path]; ok {
			clearDeadlines(w, s.log)
		}
		next.ServeHTTP(w, r)
	})
}

// clearDeadlines removes the read and write deadlines from the
// underlying connection. A zero time means "no deadline".
func clearDeadlines(w http.ResponseWriter, log *slog.Logger) {
	rc := http.NewResponseController(w)
	if err := rc.SetReadDeadline(time.Time{}); err != nil && !errors.Is(err, http.ErrNotSupported) {
		log.Warn("failed to clear read deadline for long-running request", "error", err)
	}
	if err := rc.SetWriteDeadline(time.Time{}); err != nil && !errors.Is(err, http.ErrNotSupported) {
		log.Warn("failed to clear write deadline for long-running request", "error", err)
	}
}

// wrapAuth applies the authn middleware, skipping public paths.
// Public paths are checked by exact match first, then by prefix.
// After authn sets the transport-level auth info, bridgeUserInfo
// copies it into the domain-level core.UserInfo context key so that
// infrastructure adapters can access the user identity without
// depending on the connectrpc/authn package.
func (s *Server) wrapAuth(next http.Handler) http.Handler {
	protected := s.authMiddleware.Wrap(bridgeUserInfo(next))
	if len(s.publicPaths) == 0 && len(s.publicPathPrefixes) == 0 {
		return protected
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
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
