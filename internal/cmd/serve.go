package cmd

import (
	"github.com/spf13/cobra"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/stack/v1"
	"github.com/openhdc/openhdc/internal/app"
)

func NewCmdServe(sa *app.StackApp) *cobra.Command {
	var network string
	var address string

	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv := openhdc.NewDefaultServer(
				openhdc.WithNetwork(network),
				openhdc.WithAddress(address),
			)

			grpcSrv := srv.GRPCServer()
			pb.RegisterStackServiceServer(grpcSrv, sa)

			app := openhdc.New(
				openhdc.WithServers(srv),
			)
			return app.Run()
		},
	}

	cmd.PersistentFlags().StringVar(
		&network,
		"network",
		"tcp",
		"network of grpc server",
	)

	cmd.PersistentFlags().StringVar(
		&address,
		"address",
		":0",
		"address of grpc server",
	)

	return cmd
}
