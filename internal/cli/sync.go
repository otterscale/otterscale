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
	r, err := workload.NewReader(args)
	if err != nil {
		return err
	}

	// new source clients
	srcs, err := newProcesses(ctx, r.Sources)
	if err != nil {
		return err
	}

	// new destination clients
	dsts, err := newProcesses(ctx, r.Destinations)
	if err != nil {
		return err
	}

	// new error group
	eg, gctx := errgroup.WithContext(ctx)

	// pulls
	pulls := []grpc.ServerStreamingClient[pb.Message]{}
	for _, src := range srcs {
		req := &pb.PullRequest{}
		req.SetTables(src.Tables())
		req.SetTables(src.SkipTables())

		pull, err := src.Client.Pull(gctx, req, grpc.WaitForReady(true))
		if err != nil {
			return err
		}
		pulls = append(pulls, pull)
	}

	// pushes
	pushes := []grpc.ClientStreamingClient[pb.Message, emptypb.Empty]{}
	for _, dst := range dsts {
		push, err := dst.Client.Push(gctx, grpc.WaitForReady(true))
		if err != nil {
			return err
		}
		pushes = append(pushes, push)
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

	// wait
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	// all ok
	for _, push := range pushes {
		if _, err := push.CloseAndRecv(); err != nil {
			return err
		}
	}

	// close
	for _, p := range append(srcs, dsts...) {
		if _, err := p.Client.Close(ctx, &pb.CloseRequest{}); err != nil {
			return err
		}
	}

	return nil
}
