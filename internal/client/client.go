package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/avast/retry-go/v4"
	"google.golang.org/grpc"

	"github.com/openhdc/openhdc/pkg/transport"
)

type Client struct {
	name    string
	version string
	path    string
	socket  string
	// wg     *sync.WaitGroup

	Conn *grpc.ClientConn
}

func (c *Client) Name() string {
	return c.name
}

// TODO DOWNLOAD

func (c *Client) download(_ context.Context) error {
	return nil
}

// TODO READ LOG
// TODO ERROR HANDLING

func (c *Client) start(ctx context.Context) error {
	args := []string{"serve", "--address", c.socket}
	cmd := exec.CommandContext(ctx, c.path, args...)
	cmd.SysProcAttr = sysProcAttr()
	return cmd.Start()
}

func (c *Client) exec(ctx context.Context) error {
	return retry.Do(
		func() error {
			if err := c.start(ctx); err != nil {
				return err
			}
			// c.wg.Add(1)
			conn, err := transport.NewClient(
				transport.WithEndpoint(fmt.Sprintf("unix://%s", c.socket)),
				transport.WithConnector(),
			)
			if err != nil {
				return err
			}
			c.Conn = conn
			return nil
		},
		retry.Attempts(3),
		retry.Delay(time.Second),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			slog.Debug("failed to start connector", "error", err)
		}),
	)
}

func (c *Client) Terminate() error {
	_ = os.RemoveAll(c.socket)

	return nil
}

func New(ctx context.Context, opts ...Option) (*Client, error) {
	c := &Client{
		// wg: &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt(c)
	}
	f, err := os.CreateTemp("", "openhdc-*.sock")
	if err != nil {
		return nil, err
	}
	c.socket = f.Name()
	if err := c.download(ctx); err != nil {
		return nil, err
	}
	if err := c.exec(ctx); err != nil {
		return nil, err
	}
	return c, nil
}
