package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
	"github.com/openhdc/openhdc/internal/process"
)

func NewCmdSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE:    sync,
	}
	return cmd
}

func newProcesses(ctx context.Context, wls []*workload.Workload) ([]*process.Process, error) {
	ps := []*process.Process{}
	for _, wl := range wls {
		md := wl.GetMetadata()
		if md == nil {
			return nil, fmt.Errorf("metadata is empty: %s", wl.GetInternal().GetFilePath())
		}
		p := process.New(ctx,
			process.WithName(md.GetName()),
			process.WithVersion(md.GetVersion()),
			process.WithPath(md.GetPath()),
			process.WithSpec(wl.GetSpec()),
		)
		if err := p.Download(ctx); err != nil {
			return nil, err
		}
		if err := p.Start(ctx); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func sync(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// convert args
	reader, err := workload.NewReader(args)
	if err != nil {
		return err
	}

	// new source clients
	sources, err := newProcesses(ctx, reader.Sources)
	if err != nil {
		return err
	}

	// new destination clients
	destinations, err := newProcesses(ctx, reader.Destinations)
	if err != nil {
		return err
	}

	eg, _ := errgroup.WithContext(ctx)

	// pushes
	pushes := []grpc.ClientStreamingClient[pb.Message, emptypb.Empty]{}
	for _, destination := range destinations {
		push, err := newPushClient(ctx, destination)
		if err != nil {
			return err
		}
		pushes = append(pushes, push)
	}

	// pulls
	pulls := []grpc.ServerStreamingClient[pb.Message]{}
	for _, source := range sources {
		pull, err := newPullClient(ctx, source)
		if err != nil {
			return err
		}
		pulls = append(pulls, pull)
	}

	// start sync
	for _, pull := range pulls {
		eg.Go(func() error {
			for {
				msg, err := pull.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					return err
				}
				for _, push := range pushes {
					if err := push.Send(msg); err != nil {
						if errors.Is(err, io.EOF) {
							if _, err := push.CloseAndRecv(); err != nil {
								return err
							}
						}
						return err
					}
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	// all ok
	for _, push := range pushes {
		if _, err := push.CloseAndRecv(); err != nil {
			return err
		}
	}
	return nil
}

func newPullClient(ctx context.Context, p *process.Process) (grpc.ServerStreamingClient[pb.Message], error) {
	req := &pb.PullRequest{}
	req.SetTables(p.Tables())
	req.SetTables(p.SkipTables())

	c := pb.NewConnectorClient(p.Conn)
	return c.Pull(ctx, req, grpc.WaitForReady(true))
}

func newPushClient(ctx context.Context, p *process.Process) (grpc.ClientStreamingClient[pb.Message, emptypb.Empty], error) {
	c := pb.NewConnectorClient(p.Conn)
	return c.Push(ctx, grpc.WaitForReady(true))
}
