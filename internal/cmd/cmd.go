package cmd

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"
)

const (
	containerEnvVar            = "OTTERSCALE_CONTAINER"
	defaultContainerAddress    = ":8299"
	defaultContainerConfigPath = "/etc/app/otterscale.yaml"
)

func startHTTPServer(address string, handler http.Handler) error {
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	srv := &http.Server{
		Addr:              address,
		Handler:           cors.AllowAll().Handler(handler),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
		Protocols:         protocols,
	}

	listener, err := net.Listen("tcp", address) //nolint:noctx // context not needed for Listen
	if err != nil {
		return err
	}

	slog.Info("Server starting on", "address", listener.Addr().String())
	return srv.Serve(listener)
}
