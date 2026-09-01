package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// NewJoinCommand groups what an operator hands to an agent so it can join a
// cluster to this server: the token that authorizes the cluster, and the CA
// that lets the agent verify the server it registers against.
//
// Both are derived from material already mounted in the server, so these run
// inside it:
//
//	kubectl exec deploy/otterscale-server -- \
//	    /otterscale join token --cluster prod
//
// Holding the root secret is what authorizes issuing a token: no separate admin
// credential, and nothing has to be running.
func NewJoinCommand(conf *config.Config) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Print what an agent needs to join a cluster to this server",
		Args:  cobra.NoArgs,
	}

	token, err := newJoinTokenCommand(conf)
	if err != nil {
		return nil, err
	}
	ca, err := newJoinCACommand(conf)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(token, ca)

	return cmd, nil
}

func newJoinTokenCommand(conf *config.Config) (*cobra.Command, error) {
	var cluster string

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Print the join token an agent needs to register a cluster",
		Long: "Print the join token an agent needs to register a cluster.\n\n" +
			"The token is derived from the server's join secret and the cluster " +
			"name, so it authorizes that cluster only: an agent holding it cannot " +
			"register under any other name. Tokens do not expire; rotating the " +
			"join secret invalidates all of them.",
		Example: "otterscale join token --cluster prod",
		Args:    cobra.NoArgs,
		// The token is the only output, so it pipes straight into a helm
		// invocation; usage noise would corrupt that.
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cluster == "" {
				return errors.New("--cluster is required")
			}
			if err := core.ValidateClusterName(cluster); err != nil {
				return err
			}

			secret, err := conf.ServerJoinSecret()
			if err != nil {
				return err
			}
			authority, err := core.NewJoinAuthority(secret)
			if err != nil {
				return errors.New(
					"join secret is not configured; " +
						"set --join-secret, --join-secret-file or OTTERSCALE_SERVER_JOIN_SECRET",
				)
			}

			fmt.Fprintln(cmd.OutOrStdout(), authority.Token(cluster))
			return nil
		},
	}

	cmd.Flags().StringVar(&cluster, "cluster", "", "Cluster the token authorizes")

	// Only the join flags matter for deriving a token.
	if err := conf.BindFlags(cmd.Flags(), config.JoinOptions); err != nil {
		return nil, err
	}

	return cmd, nil
}

// newJoinCACommand copies the configured CA to stdout verbatim, so the output
// is a usable PEM file rather than something a caller has to reformat.
//
// An unconfigured CA is an error rather than empty output: an empty ca.crt
// would be written out and fail later at the agent, where the cause is much
// harder to see.
func newJoinCACommand(conf *config.Config) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "ca",
		Short: "Print the CA an agent must trust to reach this server",
		Long: "Print the CA an agent must trust to reach this server.\n\n" +
			"An agent verifies this server with its image's system roots, so a " +
			"privately signed certificate has to be handed over out of band before " +
			"registration can succeed. Nothing is needed when the certificate " +
			"chains to a public CA, and this command says so.",
		Example:      "otterscale join ca > ca.crt",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			path := conf.ServerTrustedCAFile()
			if path == "" {
				return errors.New(
					"no trusted CA is configured; either this server's certificate chains to a " +
						"public CA and the agent needs nothing, or trustedCA is unset in the chart",
				)
			}

			pem, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read the trusted CA from %q: %w", path, err)
			}

			_, err = cmd.OutOrStdout().Write(pem)
			return err
		},
	}

	if err := conf.BindFlags(cmd.Flags(), config.TrustedCAOptions); err != nil {
		return nil, err
	}

	return cmd, nil
}
