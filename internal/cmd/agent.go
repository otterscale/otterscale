package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/cmd/agent"
	"github.com/otterscale/otterscale/internal/config"
)

// AgentInjector is a Wire-generated factory that creates a fully
// wired Agent together with a cleanup function.
type AgentInjector func() (*agent.Agent, func(), error)

// NewAgentCommand returns the "agent" Cobra subcommand. The injector
// is called lazily inside RunE so that expensive initialisation
// (loading kubeconfig, etc.) only happens when the command actually
// executes.
func NewAgentCommand(conf *config.Config, newAgent AgentInjector) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "agent",
		Short:   "Start agent that connects to server and executes requests in-cluster",
		Example: "otterscale agent --cluster=default --server-url=https://api.otterscale.io --tunnel-server-url=https://tunnel.otterscale.io",
		RunE: func(cmd *cobra.Command, _ []string) error {
			agt, cleanup, err := newAgent()
			if err != nil {
				return fmt.Errorf("failed to initialize agent: %w", err)
			}
			defer cleanup()

			cfg := &agent.Config{
				Cluster:         conf.AgentCluster(),
				ServerURL:       conf.AgentServerURL(),
				TunnelServerURL: conf.AgentTunnelServerURL(),
				Bootstrap:       conf.AgentBootstrap(),
			}

			return agt.Run(cmd.Context(), cfg)
		},
	}

	if err := conf.BindFlags(cmd.Flags(), config.AgentOptions); err != nil {
		return nil, err
	}

	return cmd, nil
}
