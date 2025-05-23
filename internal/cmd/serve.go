package cmd

import (
	"log"
	"net"
	"net/http"
	"time"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	applicationv1 "github.com/openhdc/otterscale/api/application/v1/pbconnect"
	configurationv1 "github.com/openhdc/otterscale/api/configuration/v1/pbconnect"
	facilityv1 "github.com/openhdc/otterscale/api/facility/v1/pbconnect"
	generalv1 "github.com/openhdc/otterscale/api/general/v1/pbconnect"
	machinev1 "github.com/openhdc/otterscale/api/machine/v1/pbconnect"
	networkv1 "github.com/openhdc/otterscale/api/network/v1/pbconnect"
	scopev1 "github.com/openhdc/otterscale/api/scope/v1/pbconnect"
	tagv1 "github.com/openhdc/otterscale/api/tag/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/app"
)

func NewCmdServe(app *app.ApplicationService, config *app.ConfigurationService, facility *app.FacilityService, general *app.GeneralService, machine *app.MachineService, network *app.NetworkService, scope *app.ScopeService, tag *app.TagService) *cobra.Command {
	var address, configPath string

	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			mux := http.NewServeMux()
			mux.Handle(applicationv1.NewApplicationServiceHandler(app))
			mux.Handle(configurationv1.NewConfigurationServiceHandler(config))
			mux.Handle(facilityv1.NewFacilityServiceHandler(facility))
			mux.Handle(generalv1.NewGeneralServiceHandler(general))
			mux.Handle(machinev1.NewMachineServiceHandler(machine))
			mux.Handle(networkv1.NewNetworkServiceHandler(network))
			mux.Handle(scopev1.NewScopeServiceHandler(scope))
			mux.Handle(tagv1.NewTagServiceHandler(tag))

			services := []string{
				applicationv1.ApplicationServiceName,
				configurationv1.ConfigurationServiceName,
				facilityv1.FacilityServiceName,
				generalv1.GeneralServiceName,
				machinev1.MachineServiceName,
				networkv1.NetworkServiceName,
				scopev1.ScopeServiceName,
				tagv1.TagServiceName,
			}

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

			listener, err := net.Listen("tcp", address)
			if err != nil {
				return err
			}

			log.Printf("Server starting on %s\n", listener.Addr().String())
			return srv.Serve(listener)
		},
	}

	cmd.PersistentFlags().StringVar(
		&address,
		"address",
		":0",
		"address of grpc server",
	)

	cmd.PersistentFlags().StringVar(
		&configPath,
		"config",
		"./otterscale.yaml",
		"path of config",
	)

	return cmd
}
