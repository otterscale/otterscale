package process

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
	"github.com/openhdc/openhdc/api/property/v1"
)

// TODO: BETTER
const maxMsgSize = 100 * 1024 * 1024

type Process struct {
	opts options

	Client pb.ConnectorClient
}

func New(ctx context.Context, opts ...Option) *Process {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Process{
		opts: o,
	}
}

func (c *Process) Name() string {
	return c.opts.name
}

func (c *Process) Tables() []string {
	// TODO: FROM WORKLOAD
	return []string{}
}

func (c *Process) SkipTables() []string {
	// TODO: FROM WORKLOAD
	return []string{}
}

// TODO DOWNLOAD

func (c *Process) Download(_ context.Context) error {
	return nil
}

// TODO READ LOG
// TODO ERROR HANDLING

func (c *Process) Start(ctx context.Context) error {
	// get address
	address, err := freeAddress()
	if err != nil {
		return err
	}

	// new grpc connection
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		))
	if err != nil {
		return err
	}

	// new grpc client
	c.Client = pb.NewConnectorClient(conn)

	// prepare arguments
	args := []string{}
	args = append(args, addressToArgs(address)...)
	args = append(args, syncModeToArgs(c.opts.syncMode)...)
	args = append(args, fieldsToArgs(c.opts.spec.GetFields())...)

	cmd := exec.CommandContext(ctx, c.opts.path, args...) //nolint:gosec
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = sysProcAttr()
	return cmd.Start()
}

func (c *Process) Terminate() error {
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

func addressToArgs(address string) []string {
	return []string{"--address", address}
}

func syncModeToArgs(sm property.SyncMode) []string {
	if sm == property.SyncMode_sync_mode_unspecified {
		return nil
	}
	return []string{"--sync_mode", sm.String()}
}

func valueToArgs(v *structpb.Value) []string {
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
			args = append(args, valueToArgs(v)...)
		}
	}
	return args
}

func fieldsToArgs(m map[string]*structpb.Value) []string {
	args := []string{}
	for key, val := range m {
		v := valueToArgs(val)
		if len(v) > 0 {
			args = append(args, "--"+key, strings.Join(v, ","))
		}
	}
	return args
}
