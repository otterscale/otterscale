// Package main is the entry point for the otterscale binary. It
// supports two subcommands:
//
//   - server: runs the control-plane (gRPC API + tunnel listener)
//   - agent:  runs inside a Kubernetes cluster and reverse-proxies
//     API requests through the tunnel
//
// Dependencies are assembled via Google Wire; see wire.go.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/cmd"
	"github.com/otterscale/otterscale/internal/cmd/agent"
	"github.com/otterscale/otterscale/internal/cmd/server"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/pki"
)

// version is injected at build time via -ldflags "-X main.version=v1.2.3".
var version = "devel"

func main() {
	if err := run(); err != nil {
		// Cobra runs with SilenceErrors, so the printing happens here.
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	rootCmd, cleanup, err := wireCmd()
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return rootCmd.ExecuteContext(ctx)
}

// newCmd is a Wire provider for the root command. Closures capture the version
// on its way to the injectors, leaving the Injector signatures unchanged.
func newCmd(conf *config.Config) (*cobra.Command, error) {
	c := &cobra.Command{
		Use:           "otterscale",
		Short:         "OtterScale: A unified platform for simplified compute, storage, and networking.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	v := core.Version(version)

	serverCmd, err := cmd.NewServerCommand(conf, func() (*server.Server, func(), error) {
		return wireServer(v, conf)
	})
	if err != nil {
		return nil, err
	}

	agentCmd, err := cmd.NewAgentCommand(conf, func() (*agent.Agent, func(), error) {
		return wireAgent(v, conf)
	})
	if err != nil {
		return nil, err
	}

	enrolmentCmd, err := cmd.NewEnrolmentTokenCommand(conf)
	if err != nil {
		return nil, err
	}

	c.AddCommand(serverCmd, agentCmd, enrolmentCmd)

	return c, nil
}

// provideCA generates a fresh ephemeral CA per server start. It is never
// persisted; agents re-register through the public Register API after a
// restart.
func provideCA() (*pki.CA, error) {
	return pki.NewCA()
}
