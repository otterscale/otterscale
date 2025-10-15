package cmd

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	containerEnvVar            = "OTTERSCALE_CONTAINER"
	defaultContainerAddress    = ":8299"
	defaultContainerConfigPath = "/etc/app/otterscale.yaml"
)

func startHTTPServer(address string, handler http.Handler) error {
	srv := &http.Server{
		Addr: address,
		Handler: h2c.NewHandler(
			cors.AllowAll().Handler(handler),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	listener, err := net.Listen("tcp", address) //nolint:noctx // context not needed for Listen
	if err != nil {
		return err
	}

	slog.Info("Server starting on", "address", listener.Addr().String())
	return srv.Serve(listener)
}
