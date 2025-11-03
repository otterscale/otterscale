package cmd

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/mux"
)

func NewJWKSProxy(jwksProxy *mux.JWKSProxy) *cobra.Command {
	var (
		address, jwksURL                       string
		maxRefreshInterval, minRefreshInterval time.Duration
	)

	cmd := &cobra.Command{
		Use:     "jwks-proxy",
		Short:   "Start the JWKS proxy server",
		Long:    "Starts a proxy server that caches and serves JWKS (JSON Web Key Set) from an external identity provider. The proxy periodically refreshes the keys to ensure they are up-to-date while reducing load on the external provider.",
		Example: "otterscale jwks-proxy --address=:8299 --jwks-url=https://default.idp.com/.well-known/jwks.json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if os.Getenv(containerEnvVar) != "" {
				address = defaultContainerAddress
				slog.Info("Container environment detected, using default configuration", "address", address)
			}

			cache, err := newJWKSCache(cmd.Context(), jwksURL, maxRefreshInterval, minRefreshInterval)
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

	cmd.Flags().DurationVar(
		&maxRefreshInterval,
		"max-refresh-interval",
		24*time.Hour*7, //nolint:mnd // reasonable default
		"Interval to check and refresh the external JWKS keys",
	)

	cmd.Flags().DurationVar(
		&minRefreshInterval,
		"min-refresh-interval",
		15*time.Minute, //nolint:mnd // reasonable default
		"Minimum interval to check and refresh the external JWKS keys",
	)

	return cmd
}

func newJWKSCache(ctx context.Context, url string, maxRefreshInterval, minRefreshInterval time.Duration) (*jwk.Cache, error) {
	cache, err := jwk.NewCache(
		ctx,
		httprc.NewClient(
			httprc.WithTraceSink(tracesink.NewSlog(slog.New(slog.NewJSONHandler(os.Stderr, nil)))),
		))
	if err != nil {
		return nil, err
	}
	if err := cache.Register(
		ctx,
		url,
		jwk.WithMaxInterval(maxRefreshInterval),
		jwk.WithMinInterval(minRefreshInterval),
	); err != nil {
		return nil, err
	}
	if _, err := cache.Refresh(ctx, url); err != nil {
		return nil, err
	}
	return cache, nil
}
