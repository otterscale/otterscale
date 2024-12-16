package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

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
		md := wl.Metadata
		if md == nil {
			return nil, fmt.Errorf("metadata is empty: %s", wl.Internal.FilePath)
		}
		spec := wl.Spec
		if spec == nil {
			return nil, fmt.Errorf("spec is empty: %s", wl.Internal.FilePath)
		}
		p := process.New(ctx,
			process.WithName(md.Name),
			process.WithVersion(md.Version),
			process.WithPath(md.Path),
		)
		if err := p.Download(ctx); err != nil {
			return nil, err
		}
		if err := p.Start(ctx, spec); err != nil {
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

	// new transformer clients
	transformers, err := newProcesses(ctx, reader.Transformers)
	if err != nil {
		return err
	}

	// start sync
	for _, source := range sources {
		fmt.Println(source.Name())

		srcClient := pb.NewConnectorClient(source.Conn)

		req := &pb.PullRequest{
			Tables: []string{},
		}

		pull, err := srcClient.Pull(ctx, req, grpc.WaitForReady(true))
		if err != nil {
			return err
		}

		r, err := pull.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		for _, destination := range destinations {
			for _, transformer := range transformers {
				fmt.Println(transformer.Name())
			}

			fmt.Println(destination.Name())
			dstClient := pb.NewConnectorClient(destination.Conn)

			rs, err := dstClient.Push(ctx, grpc.WaitForReady(true))
			if err != nil {
				return err
			}

			if err := rs.Send(r); err != nil {
				return err
			}

			time.Sleep(time.Second * 5)
		}
	}

	return nil
}
