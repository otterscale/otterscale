package cmd

import (
	"context"
	"fmt"

	"github.com/openhdc/openhdc/internal/connector"
	"github.com/openhdc/openhdc/internal/workload"
	"github.com/spf13/cobra"
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

func sourcesToConnectors(ctx context.Context, sources []*workload.Source) ([]*connector.Connector, error) {
	connectors := []*connector.Connector{}
	for _, source := range sources {
		opts := []connector.Option{
			connector.WithName(source.Name),
			connector.WithVersion(source.Version),
			connector.WithPath(source.Path),
		}
		connector, err := connector.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
		connectors = append(connectors, connector)
	}
	return connectors, nil
}

func destinationsToConnectors(ctx context.Context, destinations []*workload.Destination) ([]*connector.Connector, error) {
	connectors := []*connector.Connector{}
	for _, destination := range destinations {
		opts := []connector.Option{
			connector.WithName(destination.Name),
			connector.WithVersion(destination.Version),
			connector.WithPath(destination.Path),
		}
		connector, err := connector.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
		connectors = append(connectors, connector)
	}
	return connectors, nil
}

func transformersToConnectors(ctx context.Context, transformers []*workload.Transformer) ([]*connector.Connector, error) {
	connectors := []*connector.Connector{}
	for _, transformer := range transformers {
		opts := []connector.Option{
			connector.WithName(transformer.Name),
			connector.WithVersion(transformer.Version),
			connector.WithPath(transformer.Path),
		}
		connector, err := connector.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
		connectors = append(connectors, connector)
	}
	return connectors, nil
}

func sync(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// convert args
	reader, err := workload.NewReader(args)
	if err != nil {
		return err
	}

	// new source connectors
	sources, err := sourcesToConnectors(ctx, reader.Sources)
	if err != nil {
		return err
	}

	// new destination connectors
	destinations, err := destinationsToConnectors(ctx, reader.Destinations)
	if err != nil {
		return err
	}

	// new transformer connectors
	transformers, err := transformersToConnectors(ctx, reader.Transformers)
	if err != nil {
		return err
	}

	// start sync
	for _, source := range sources {
		fmt.Println(source.Name())
		for _, destination := range destinations {
			fmt.Println(destination.Name())
			for _, transformer := range transformers {
				fmt.Println(transformer.Name())
			}
		}
	}

	return nil
}
