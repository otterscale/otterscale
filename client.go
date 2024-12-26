package openhdc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
)

type Client struct {
	opts options

	pb.ConnectorClient
}

type options struct {
	path        string
	sync        *workload.Sync
	spec        *structpb.Struct
	grpcOptions []grpc.DialOption
}

var defaultOptions = options{
	grpcOptions: []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultMaxMessageSize),
			grpc.MaxCallSendMsgSize(defaultMaxMessageSize),
		),
	},
}

type Option interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

var _ Option = (*funcOption)(nil)

func (fro *funcOption) apply(ro *options) {
	fro.f(ro)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithPath(p string) Option {
	return newFuncOption(func(o *options) {
		o.path = p
	})
}

func WithSync(s *workload.Sync) Option {
	return newFuncOption(func(o *options) {
		o.sync = s
	})
}

func WithSpec(s *structpb.Struct) Option {
	return newFuncOption(func(o *options) {
		o.spec = s
	})
}

func NewClient(opt ...Option) *Client {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	c := &Client{
		opts: opts,
	}
	return c
}

func (c *Client) Sync() *workload.Sync {
	return c.opts.sync
}

// TODO: DOWNLOAD
func (c *Client) download(_ context.Context) error {
	return nil
}

func (c *Client) Start(ctx context.Context) error {
	if err := c.download(ctx); err != nil {
		return err
	}

	target, err := c.freePort()
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(target, c.opts.grpcOptions...)
	if err != nil {
		return err
	}
	c.ConnectorClient = pb.NewConnectorClient(conn)

	args := []string{"--address", target}
	args = append(args, c.specToArgs(c.opts.spec)...)

	cmd := exec.CommandContext(ctx, c.opts.path, args...) //nolint:gosec
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = sysProcAttr()
	return cmd.Start()
}

// TODO: TERMINATE
func (c *Client) Terminate() error {
	return nil
}

func (c *Client) freePort() (string, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	defer lis.Close()
	return lis.Addr().String(), nil
}

func (c *Client) specToArgs(s *structpb.Struct) []string {
	args := []string{}
	for key, val := range s.GetFields() {
		v := c.valueToArgs(val)
		if len(v) > 0 {
			args = append(args, "--"+key, strings.Join(v, ","))
		}
	}
	return args
}

func (c *Client) valueToArgs(v *structpb.Value) []string {
	args := []string{}
	switch v.GetKind().(type) {
	case *structpb.Value_NullValue:
		break
	case *structpb.Value_NumberValue:
		args = append(args, fmt.Sprintf("%v", v.GetNumberValue()))
	case *structpb.Value_StringValue:
		args = append(args, v.GetStringValue())
	case *structpb.Value_BoolValue:
		args = append(args, fmt.Sprintf("%v", v.GetBoolValue()))
	case *structpb.Value_StructValue:
		// not support
	case *structpb.Value_ListValue:
		for _, v := range v.GetListValue().GetValues() {
			args = append(args, c.valueToArgs(v)...)
		}
	}
	return args
}
