package cmd

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/mux"
)

func NewJWKSProxy(jwksProxy *mux.JWKSProxy) *cobra.Command {
	var (
		address, jwksURL string
		refreshInterval  time.Duration
	)

	cmd := &cobra.Command{
		Use:     "jwks-proxy",
		Short:   "Start the JWKS proxy server",
		Long:    "Start the OtterScale API server that provides gRPC and HTTP endpoints for JWKS proxy service",
		Example: "otterscale jwks-proxy --address=:8299",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if os.Getenv(containerEnvVar) != "" {
				address = defaultContainerAddress
				slog.Info("Container environment detected, using default configuration", "address", address)
			}

			cache, err := newJWKSCache(cmd.Context(), jwksURL)
			if err != nil {
				return err
			}
			jwksProxy.SetCache(cache)
			jwksProxy.SetURL(jwksURL)

			return startHTTPServer(address, jwksProxy)
		},
	}

	cmd.Flags().StringVarP(
		&address,
		"address",
		"a",
		":0",
		"Address for server to listen on",
	)

	cmd.Flags().StringVarP(
		&jwksURL,
		"jwks-url",
		"u",
		"https://default.idp.com/.well-known/jwks.json",
		"External JWKS endpoint URL",
	)

	cmd.Flags().DurationVarP(
		&refreshInterval,
		"refresh-interval",
		"i",
		5*time.Minute,
		"Interval to check and refresh the external JWKS keys",
	)

	return cmd
}

// TODO: Check Interval
func newJWKSCache(ctx context.Context, url string) (*jwk.Cache, error) {
	cache, err := jwk.NewCache(ctx, httprc.NewClient())
	if err != nil {
		return nil, err
	}
	if err := cache.Register(ctx, url); err != nil {
		return nil, err
	}
	if _, err := cache.Refresh(ctx, url); err != nil {
		return nil, err
	}
	return cache, nil
}
