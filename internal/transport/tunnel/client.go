package tunnel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	chclient "github.com/jpillora/chisel/client"
)

// Sentinel errors for well-known failure modes.
var (
	ErrLocalPortRequired = errors.New("tunnel: local port is required")
	ErrRegisterRequired  = errors.New("tunnel: register function is required")
)

// RegisterResult holds the mTLS credentials and tunnel endpoint
// returned by a successful registration.
type RegisterResult struct {
	// Endpoint is the tunnel server's allocated address (host:port).
	Endpoint string
	// Auth is the chisel auth string ("user:password") derived from
	// the signed certificate.
	Auth string
	// CACertPEM is the PEM-encoded CA certificate for verifying the
	// tunnel server.
	CACertPEM []byte
	// CertPEM is the PEM-encoded client certificate signed by the CA.
	CertPEM []byte
	// KeyPEM is the PEM-encoded private key corresponding to the
	// client certificate.
	KeyPEM []byte
}

// RegisterFunc registers an agent and returns mTLS credentials.
type RegisterFunc func(ctx context.Context, serverURL, cluster string) (*RegisterResult, error)

// ClientOption configures a Client.
type ClientOption func(*Client)

// Client manages a reverse tunnel connection with automatic
// registration, reconnection, and exponential backoff. It uses mTLS
// for tunnel authentication.
type Client struct {
	mu               sync.Mutex       // protects inner and certDir
	inner            *chclient.Client // owned lifecycle, not exported
	certDir          string           // temp directory for TLS cert files

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

// WithCluster configures the cluster name used for registration.
func WithCluster(cluster string) ClientOption {
	return func(c *Client) { c.cluster = cluster }
}

// WithServerURL configures the fleet server URL for registration.
func WithServerURL(serverURL string) ClientOption {
	return func(c *Client) { c.serverURL = serverURL }
}

// WithTunnelServerURL configures the chisel tunnel server URL.
func WithTunnelServerURL(tunnelServerURL string) ClientOption {
	return func(c *Client) { c.tunnelServerURL = tunnelServerURL }
}

// WithLocalPort configures the local port to expose through the tunnel.
func WithLocalPort(localPort int) ClientOption {
	return func(c *Client) { c.localPort = localPort }
}

// WithKeepAlive configures the keep-alive interval for the tunnel.
func WithKeepAlive(keepAlive time.Duration) ClientOption {
	return func(c *Client) { c.keepAlive = keepAlive }
}

// WithMaxRetryCount configures chisel's internal maximum retry count.
func WithMaxRetryCount(maxRetryCount int) ClientOption {
	return func(c *Client) { c.maxRetryCount = maxRetryCount }
}

// WithMaxRetryInterval configures chisel's internal maximum retry interval.
func WithMaxRetryInterval(maxRetryInterval time.Duration) ClientOption {
	return func(c *Client) { c.maxRetryInterval = maxRetryInterval }
}

// WithBaseRetryDelay configures the initial delay for the outer reconnect backoff.
func WithBaseRetryDelay(baseRetryDelay time.Duration) ClientOption {
	return func(c *Client) { c.baseRetryDelay = baseRetryDelay }
}

// WithMaxRetryDelay configures the maximum delay for the outer reconnect backoff.
func WithMaxRetryDelay(maxRetryDelay time.Duration) ClientOption {
	return func(c *Client) { c.maxRetryDelay = maxRetryDelay }
}

// WithRegister configures the function used to register with the fleet server.
func WithRegister(register RegisterFunc) ClientOption {
	return func(c *Client) { c.register = register }
}

// WithLogger configures a structured logger. Defaults to slog.Default
// with "component" and "cluster" attributes.
func WithLogger(log *slog.Logger) ClientOption {
	return func(c *Client) { c.log = log }
}

// NewClient creates a tunnel client. It validates required fields
// but does not perform any I/O.
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		cluster:          "default",
		serverURL:        "http://127.0.0.1:8299",
		tunnelServerURL:  "https://127.0.0.1:8300",
		keepAlive:        30 * time.Second,
		maxRetryCount:    3,
		maxRetryInterval: 10 * time.Second,
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

// Start runs the tunnel client loop. It blocks until ctx is cancelled,
// automatically re-registering and reconnecting on failures with
// exponential backoff.
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
		bo.Reset()
		c.mu.Lock()
		c.inner = inner
		c.mu.Unlock()

		err = c.runSession(ctx, inner)
		if ctx.Err() != nil {
			return nil
		}
		if err == nil || isAuthErr(err) {
			if err != nil {
				c.log.Warn("authentication failed, re-registering", "error", err)
			} else {
				c.log.Warn("session ended, re-registering")
			}
			bo.Reset()
			continue
		}

		c.log.Warn("connection lost, retrying", "error", err, "retry_in", bo.current)
		if !sleepCtx(ctx, bo.Next()) {
			return nil
		}
	}
}

// Stop gracefully shuts down the tunnel client and cleans up temp files.
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

// dial registers with the fleet server, writes mTLS credentials to
// temp files, and creates a new chisel client configured for mTLS.
func (c *Client) dial(ctx context.Context) (*chclient.Client, error) {
	result, err := c.register(ctx, c.serverURL, c.cluster)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	c.log.Info("registered", "endpoint", result.Endpoint)

	// Write mTLS credentials to a temp directory.
	dir, err := os.MkdirTemp("", "otterscale-tls-*")
	if err != nil {
		return nil, fmt.Errorf("create cert dir: %w", err)
	}

	// Atomically swap the cert directory under a single lock to
	// avoid a TOCTOU race with Stop().
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

	if err := os.WriteFile(caFile, result.CACertPEM, 0600); err != nil {
		return nil, fmt.Errorf("write CA cert: %w", err)
	}
	if err := os.WriteFile(certFile, result.CertPEM, 0600); err != nil {
		return nil, fmt.Errorf("write client cert: %w", err)
	}
	if err := os.WriteFile(keyFile, result.KeyPEM, 0600); err != nil {
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
	})
}

// runSession starts the inner chisel client and waits for it to finish.
// It always closes the inner client before returning.
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

