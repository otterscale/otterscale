package tunnel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	chclient "github.com/jpillora/chisel/client"
)

// healthySession is how long a tunnel session must last before it
// counts as a successful connection, resetting the reconnect backoff.
const healthySession = 30 * time.Second

var (
	ErrLocalPortRequired = errors.New("tunnel: local port is required")
	ErrRegisterRequired  = errors.New("tunnel: register function is required")
)

// RegisterResult holds what a successful registration returns.
type RegisterResult struct {
	// Endpoint is the tunnel server's allocated address (host:port).
	Endpoint string
	// Auth is the chisel auth string ("user:password").
	Auth string
	// CACertPEM verifies the tunnel server.
	CACertPEM []byte
	// CertPEM is the client certificate signed by the CA.
	CertPEM []byte
	// KeyPEM corresponds to CertPEM.
	KeyPEM []byte
}

// RegisterFunc registers an agent and returns mTLS credentials.
type RegisterFunc func(ctx context.Context, serverURL, cluster string) (*RegisterResult, error)

type ClientOption func(*Client)

// Client manages an mTLS reverse tunnel with automatic registration,
// reconnection, and exponential backoff.
type Client struct {
	mu      sync.Mutex       // protects inner and certDir
	inner   *chclient.Client // owned lifecycle, not exported
	certDir string           // temp directory for TLS cert files

	cluster          string
	serverURL        string
	tunnelServerURL  string
	localPort        int
	keepAlive        time.Duration
	maxRetryCount    int
	maxRetryInterval time.Duration
	baseRetryDelay   time.Duration
	maxRetryDelay    time.Duration
	register         RegisterFunc
	log              *slog.Logger
}

func WithCluster(cluster string) ClientOption {
	return func(c *Client) { c.cluster = cluster }
}

// WithServerURL sets the link server used for registration.
func WithServerURL(serverURL string) ClientOption {
	return func(c *Client) { c.serverURL = serverURL }
}

func WithTunnelServerURL(tunnelServerURL string) ClientOption {
	return func(c *Client) { c.tunnelServerURL = tunnelServerURL }
}

// WithLocalPort sets the port exposed through the tunnel.
func WithLocalPort(localPort int) ClientOption {
	return func(c *Client) { c.localPort = localPort }
}

func WithKeepAlive(keepAlive time.Duration) ClientOption {
	return func(c *Client) { c.keepAlive = keepAlive }
}

// WithMaxRetryCount bounds chisel's own retries.
func WithMaxRetryCount(maxRetryCount int) ClientOption {
	return func(c *Client) { c.maxRetryCount = maxRetryCount }
}

// WithMaxRetryInterval bounds chisel's own retry interval.
func WithMaxRetryInterval(maxRetryInterval time.Duration) ClientOption {
	return func(c *Client) { c.maxRetryInterval = maxRetryInterval }
}

// WithBaseRetryDelay sets the outer reconnect backoff's initial delay.
func WithBaseRetryDelay(baseRetryDelay time.Duration) ClientOption {
	return func(c *Client) { c.baseRetryDelay = baseRetryDelay }
}

// WithMaxRetryDelay caps the outer reconnect backoff.
func WithMaxRetryDelay(maxRetryDelay time.Duration) ClientOption {
	return func(c *Client) { c.maxRetryDelay = maxRetryDelay }
}

func WithRegister(register RegisterFunc) ClientOption {
	return func(c *Client) { c.register = register }
}

// WithLogger defaults to slog.Default with "component" and "cluster" attributes.
func WithLogger(log *slog.Logger) ClientOption {
	return func(c *Client) { c.log = log }
}

// NewClient validates required fields but performs no I/O.
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		cluster:          "default",
		serverURL:        "http://127.0.0.1:8299",
		tunnelServerURL:  "https://127.0.0.1:8300",
		keepAlive:        15 * time.Second,
		maxRetryCount:    3,
		maxRetryInterval: 5 * time.Second,
		baseRetryDelay:   1 * time.Second,
		maxRetryDelay:    30 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.localPort == 0 {
		return nil, ErrLocalPortRequired
	}
	if c.register == nil {
		return nil, ErrRegisterRequired
	}
	if c.log == nil {
		c.log = slog.Default().With("component", "tunnel-client", "cluster", c.cluster)
	}

	return c, nil
}

// Start blocks until ctx is canceled, re-registering and reconnecting on
// failure with exponential backoff.
func (c *Client) Start(ctx context.Context) error {
	bo := newBackoff(c.baseRetryDelay, c.maxRetryDelay)

	for {
		if ctx.Err() != nil {
			return nil
		}

		inner, err := c.dial(ctx)
		if err != nil {
			c.log.Warn("registration failed, retrying", "error", err, "retry_in", bo.current)
			if !sleepCtx(ctx, bo.Next()) {
				return nil
			}
			continue
		}
		c.mu.Lock()
		c.inner = inner
		c.mu.Unlock()

		start := time.Now()
		err = c.runSession(ctx, inner)
		if ctx.Err() != nil {
			return nil
		}

		// Only a session that stayed up is evidence the server is healthy.
		// Resetting on every outcome would turn a server that accepts and
		// immediately drops connections — or rejects the credentials — into an
		// unthrottled registration loop.
		if time.Since(start) >= healthySession {
			bo.Reset()
		}

		switch {
		case err == nil:
			c.log.Warn("session ended, re-registering", "retry_in", bo.current)
		case isAuthErr(err):
			c.log.Warn("authentication failed, re-registering", "error", err, "retry_in", bo.current)
		default:
			c.log.Warn("connection lost, retrying", "error", err, "retry_in", bo.current)
		}

		if !sleepCtx(ctx, bo.Next()) {
			return nil
		}
	}
}

// Stop shuts the client down and removes its temp files.
func (c *Client) Stop(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.certDir != "" {
		if err := os.RemoveAll(c.certDir); err != nil {
			c.log.Warn("failed to remove cert dir", "error", err)
		}
		c.certDir = ""
	}
	if c.inner == nil {
		return nil
	}
	c.log.Info("shutting down")
	return c.inner.Close()
}

// dial registers with the link server, writes the mTLS credentials to temp
// files, and builds a chisel client around them.
func (c *Client) dial(ctx context.Context) (*chclient.Client, error) {
	result, err := c.register(ctx, c.serverURL, c.cluster)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	c.log.Info("registered", "endpoint", result.Endpoint)

	dir, err := os.MkdirTemp("", "otterscale-tls-*")
	if err != nil {
		return nil, fmt.Errorf("create cert dir: %w", err)
	}

	// Swap under a single lock, or Stop() races with this.
	c.mu.Lock()
	oldDir := c.certDir
	c.certDir = dir
	c.mu.Unlock()

	if oldDir != "" {
		if err := os.RemoveAll(oldDir); err != nil {
			c.log.Warn("failed to remove old cert dir", "error", err)
		}
	}

	caFile := filepath.Join(dir, "ca.pem")
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")

	const secretFilePerm = 0o600 // owner-only read/write for TLS files
	if err := os.WriteFile(caFile, result.CACertPEM, secretFilePerm); err != nil {
		return nil, fmt.Errorf("write CA cert: %w", err)
	}
	if err := os.WriteFile(certFile, result.CertPEM, secretFilePerm); err != nil {
		return nil, fmt.Errorf("write client cert: %w", err)
	}
	if err := os.WriteFile(keyFile, result.KeyPEM, secretFilePerm); err != nil {
		return nil, fmt.Errorf("write client key: %w", err)
	}

	return chclient.NewClient(&chclient.Config{
		Server: c.tunnelServerURL,
		Auth:   result.Auth,
		TLS: chclient.TLSConfig{
			CA:   caFile,
			Cert: certFile,
			Key:  keyFile,
		},
		Remotes:          []string{fmt.Sprintf("R:%s:127.0.0.1:%d", result.Endpoint, c.localPort)},
		KeepAlive:        c.keepAlive,
		MaxRetryCount:    c.maxRetryCount,
		MaxRetryInterval: c.maxRetryInterval,
		DialContext:      tcpKeepAliveDialer,
	})
}

// tcpKeepAliveDialer enables aggressive TCP keepalive. Without it a half-open
// connection — a server killed without sending FIN — goes undetected for 13–30
// minutes (Linux default tcp_retries2=15); these settings cut that to 30–45
// seconds.
func tcpKeepAliveDialer(ctx context.Context, network, addr string) (net.Conn, error) {
	d := net.Dialer{
		KeepAliveConfig: net.KeepAliveConfig{
			Enable:   true,
			Idle:     15 * time.Second,
			Interval: 5 * time.Second,
			Count:    3,
		},
	}
	return d.DialContext(ctx, network, addr)
}

// runSession waits for the inner chisel client to finish, always closing it
// before returning.
func (c *Client) runSession(ctx context.Context, inner *chclient.Client) error {
	c.log.Info("connecting", "server", c.tunnelServerURL)

	if err := inner.Start(ctx); err != nil {
		if closeErr := inner.Close(); closeErr != nil {
			c.log.Warn("failed to close inner client after start failure", "error", closeErr)
		}
		return fmt.Errorf("start: %w", err)
	}

	err := inner.Wait()
	if closeErr := inner.Close(); closeErr != nil {
		c.log.Warn("failed to close inner client", "error", closeErr)
	}
	return err
}
