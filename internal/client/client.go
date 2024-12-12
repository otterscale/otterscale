package client

import (
	"context"
	"net"
	"os"
	"os/exec"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: BETTER
const maxMsgSize = 100 * 1024 * 1024

type Client struct {
	opts options

	Conn *grpc.ClientConn
}

func New(ctx context.Context, opts ...Option) *Client {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Client{
		opts: o,
	}
}

func (c *Client) Name() string {
	return c.opts.name
}

// TODO DOWNLOAD

func (c *Client) Download(_ context.Context) error {
	return nil
}

// TODO READ LOG
// TODO ERROR HANDLING

func (c *Client) Start(ctx context.Context, spec map[string]string) error {
	address, err := freeAddress()
	if err != nil {
		return err
	}

	c.Conn, err = grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		))
	if err != nil {
		return err
	}

	args := []string{"--address", address}
	for k, v := range spec {
		args = append(args, k, v)
	}
	cmd := exec.CommandContext(ctx, c.opts.path, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = sysProcAttr()
	return cmd.Start()
}

func (c *Client) Terminate() error {
	return nil
}

func freeAddress() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer lis.Close()
	return lis.Addr().String(), nil
}
