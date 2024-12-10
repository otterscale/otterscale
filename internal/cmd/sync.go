package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/client"
	"github.com/openhdc/openhdc/internal/workload"
	"github.com/openhdc/openhdc/internal/workload/spec"
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

func sourcesToClients(ctx context.Context, sources []*spec.Source) ([]*client.Client, error) {
	clients := []*client.Client{}
	for _, source := range sources {
		client, err := toClient(ctx, source)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func destinationsToClients(ctx context.Context, destinations []*spec.Destination) ([]*client.Client, error) {
	clients := []*client.Client{}
	for _, destination := range destinations {
		client, err := toClient(ctx, destination)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func transformersToClients(ctx context.Context, transformers []*spec.Transformer) ([]*client.Client, error) {
	clients := []*client.Client{}
	for _, transformer := range transformers {
		client, err := toClient(ctx, transformer)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func toClient(ctx context.Context, s spec.Spec) (*client.Client, error) {
	md := s.GetMetadata()
	c := client.New(ctx,
		client.WithName(md.Name),
		client.WithVersion(md.Version),
		client.WithPath(md.Path),
	)
	if err := c.Download(ctx); err != nil {
		return nil, err
	}
	if err := c.Start(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

func sync(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// convert args
	reader, err := workload.NewReader(args)
	if err != nil {
		return err
	}

	// new source clients
	sources, err := sourcesToClients(ctx, reader.Sources)
	if err != nil {
		return err
	}

	// new destination clients
	destinations, err := destinationsToClients(ctx, reader.Destinations)
	if err != nil {
		return err
	}

	// new transformer clients
	transformers, err := transformersToClients(ctx, reader.Transformers)
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

		rec := r.GetRecord()

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

			req1 := &pb.PushRequest{
				Record: rec,
			}

			if err := rs.Send(req1); err != nil {
				return err
			}

			time.Sleep(time.Second * 5)
		}
	}

	return nil
}
