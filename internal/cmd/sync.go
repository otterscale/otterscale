package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
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

func sync(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// new workload
	r, err := workload.NewReader(args)
	if err != nil {
		return err
	}

	// new clients
	srcClients, err := newClients(ctx, r.Sources)
	if err != nil {
		return err
	}

	dstClients, err := newClients(ctx, r.Destinations)
	if err != nil {
		return err
	}

	// new streamings
	srcStreamings, err := newSrcStreamings(ctx, srcClients)
	if err != nil {
		return err
	}

	dstStreamings, err := newDstStreamings(ctx, srcClients)
	if err != nil {
		return err
	}

	// new error group
	eg, _ := errgroup.WithContext(ctx)

	// record start time
	startedAt := time.Now()

	// start sync
	for _, src := range srcStreamings {
		eg.Go(func() error {
			return syncOneToAll(src, dstStreamings)
		})
	}

	// wait
	if err := eg.Wait(); err != nil {
		return err
	}

	// all ok
	for _, dst := range dstStreamings {
		if _, err := dst.CloseAndRecv(); err != nil {
			return err
		}
	}

	// close
	for _, c := range append(srcClients, dstClients...) {
		if _, err := c.Close(ctx, &pb.CloseRequest{}); err != nil {
			return err
		}
	}

	slog.Info("[Sync] finished", "duration", time.Since(startedAt))

	return nil
}

func newClients(ctx context.Context, ws []*workload.Workload) ([]*openhdc.Client, error) {
	cs := []*openhdc.Client{}
	for _, w := range ws {
		m := w.GetMetadata()
		if m == nil {
			return nil, fmt.Errorf("metadata is empty: %s", w.GetInternal().GetFilePath())
		}
		c := openhdc.NewClient(
			openhdc.WithPath(m.GetPath()),
			openhdc.WithSync(w.GetSync()),
			openhdc.WithSpec(w.GetSpec()),
		)
		if err := c.Start(ctx); err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func newSrcStreamings(ctx context.Context, cs []*openhdc.Client) ([]grpc.ServerStreamingClient[pb.Message], error) {
	ss := []grpc.ServerStreamingClient[pb.Message]{}
	for _, c := range cs {
		req := &pb.PullRequest{}
		req.SetSync(c.Sync())

		s, err := c.Pull(ctx, req, grpc.WaitForReady(true))
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	return ss, nil
}

func newDstStreamings(ctx context.Context, cs []*openhdc.Client) ([]grpc.ClientStreamingClient[pb.Message, emptypb.Empty], error) {
	ss := []grpc.ClientStreamingClient[pb.Message, emptypb.Empty]{}
	for _, c := range cs {
		s, err := c.Push(ctx, grpc.WaitForReady(true))
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	return ss, nil
}

func syncOneToAll(src grpc.ServerStreamingClient[pb.Message], dsts []grpc.ClientStreamingClient[pb.Message, emptypb.Empty]) error {
	var bar *progressbar.ProgressBar
	for {
		msg, err := src.Recv()
		if errors.Is(err, io.EOF) {
			_ = bar.Finish()
			slog.Info("[Sync] read finished")
			break
		}
		if err != nil {
			return err
		}
		if msg.GetKind() == property.MessageKind_migrate {
			if bar != nil {
				_ = bar.Finish()
			}
			bar = progressbar.Default(-1, "Syncing")
		}
		_ = bar.Add(1)
		for _, dst := range dsts {
			err := dst.Send(msg)
			if errors.Is(err, io.EOF) {
				slog.Error("[Sync] write error occurred")
				if _, err := dst.CloseAndRecv(); err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
