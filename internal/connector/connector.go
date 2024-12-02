package connector

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/avast/retry-go/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	maxMsgSize = 100 * 1024 * 1024
)

type Connector struct {
	path   string
	socket string
	// wg     *sync.WaitGroup

	Conn *grpc.ClientConn
}

// TODO DOWNLOAD

func (c *Connector) download(_ context.Context) error {
	return nil
}

// TODO READ LOG
// TODO ERROR HANDLING

func (c *Connector) start(ctx context.Context) error {
	args := []string{"serve", "--address", c.socket}
	cmd := exec.CommandContext(ctx, c.path, args...)
	cmd.SysProcAttr = sysProcAttr()
	return cmd.Start()
}

func (c *Connector) connection() (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
	}
	return grpc.NewClient(fmt.Sprintf("unix://%s", c.socket), opts...)
}

func (c *Connector) exec(ctx context.Context) error {
	return retry.Do(
		func() error {
			if err := c.start(ctx); err != nil {
				return err
			}
			// c.wg.Add(1)
			conn, err := c.connection()
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

func (c *Connector) Terminate() error {
	_ = os.RemoveAll(c.socket)

	return nil
}

func New(ctx context.Context, opts ...Option) (*Connector, error) {
	f, err := os.CreateTemp("", "openhdc-*.sock")
	if err != nil {
		return nil, err
	}
	c := &Connector{
		path:   "",
		socket: f.Name(),
		// wg:     &sync.WaitGroup{},
		// removeSocket: func() error { return os.Remove(f.Name()) },
	}
	for _, opt := range opts {
		opt(c)
	}
	if err := c.download(ctx); err != nil {
		return nil, err
	}
	if err := c.exec(ctx); err != nil {
		return nil, err
	}
	return c, nil
}
