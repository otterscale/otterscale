package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/client"
	"github.com/openhdc/openhdc/internal/workload"
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

func sourcesToClients(ctx context.Context, sources []*workload.Source) ([]*client.Client, error) {
	clients := []*client.Client{}
	for _, source := range sources {
		opts := []client.Option{
			client.WithName(source.Name),
			client.WithVersion(source.Version),
			client.WithPath(source.Path),
		}
		client, err := client.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func destinationsToClients(ctx context.Context, destinations []*workload.Destination) ([]*client.Client, error) {
	clients := []*client.Client{}
	for _, destination := range destinations {
		opts := []client.Option{
			client.WithName(destination.Name),
			client.WithVersion(destination.Version),
			client.WithPath(destination.Path),
		}
		client, err := client.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func transformersToClients(ctx context.Context, transformers []*workload.Transformer) ([]*client.Client, error) {
	clients := []*client.Client{}
	for _, transformer := range transformers {
		opts := []client.Option{
			client.WithName(transformer.Name),
			client.WithVersion(transformer.Version),
			client.WithPath(transformer.Path),
		}
		client, err := client.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
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
		fmt.Printf("%+v", string(rec))

		for _, destination := range destinations {
			fmt.Println(destination.Name())
			for _, transformer := range transformers {
				fmt.Println(transformer.Name())
			}
		}
	}

	return nil
}
