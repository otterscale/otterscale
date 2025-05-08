package cmd

import (
	"log"
	"net/http"
	"time"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	nexusv1 "github.com/openhdc/openhdc/api/nexus/v1/pbconnect"
	"github.com/openhdc/openhdc/internal/app"
)

func NewCmdServe(na *app.NexusApp) *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			mux := http.NewServeMux()
			mux.Handle(nexusv1.NewNexusHandler(na))

			services := []string{nexusv1.NexusName}

			checker := grpchealth.NewStaticChecker(services...)
			mux.Handle(grpchealth.NewHandler(checker))

			reflector := grpcreflect.NewStaticReflector(services...)
			mux.Handle(grpcreflect.NewHandlerV1(reflector))
			mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

			srv := &http.Server{
				Addr: address,
				Handler: h2c.NewHandler(
					cors.AllowAll().Handler(mux),
					&http2.Server{},
				),
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       5 * time.Minute,
				WriteTimeout:      5 * time.Minute,
				MaxHeaderBytes:    8 * 1024, // 8KiB
			}

			log.Printf("Server starting on %s\n", address)
			return srv.ListenAndServe()
		},
	}

	cmd.PersistentFlags().StringVar(
		&address,
		"address",
		":0",
		"address of grpc server",
	)

	return cmd
}
