package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// NewEnrolmentTokenCommand returns the "enrolment-token" subcommand,
// which prints the token an agent must present to register the given
// cluster.
//
// The token is derived from the server's root secret, so this runs
// wherever that secret is available — most conveniently inside the
// server itself, where the secret is already mounted:
//
//	kubectl exec deploy/otterscale-server -- \
//	    /otterscale enrolment-token --cluster prod
//
// Holding the root secret is what authorizes issuing a token; there is
// no separate admin credential to manage, and nothing needs to be
// running for this to work.
func NewEnrolmentTokenCommand(conf *config.Config) (*cobra.Command, error) {
	var cluster string

	cmd := &cobra.Command{
		Use:   "enrolment-token",
		Short: "Print the enrolment token an agent needs to register a cluster",
		Long: "Print the enrolment token an agent needs to register a cluster.\n\n" +
			"The token is derived from the server's enrolment secret and the cluster " +
			"name, so it authorizes that cluster only: an agent holding it cannot " +
			"register under any other name. Tokens do not expire; rotating the " +
			"enrolment secret invalidates all of them.",
		Example: "otterscale enrolment-token --cluster prod",
		Args:    cobra.NoArgs,
		// The token is the only output, so it can be piped straight
		// into a helm invocation; usage noise would corrupt that.
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cluster == "" {
				return errors.New("--cluster is required")
			}
			if err := core.ValidateClusterName(cluster); err != nil {
				return err
			}

			secret, err := conf.ServerEnrolmentSecret()
			if err != nil {
				return err
			}
			enrolment, err := core.NewEnrolment(secret)
			if err != nil {
				return errors.New(
					"enrolment secret is not configured; " +
						"set --enrolment-secret, --enrolment-secret-file or OTTERSCALE_SERVER_ENROLMENT_SECRET",
				)
			}

			fmt.Fprintln(cmd.OutOrStdout(), enrolment.Token(cluster))
			return nil
		},
	}

	cmd.Flags().StringVar(&cluster, "cluster", "", "Cluster the token authorizes")

	// Only the enrolment flags are bound: the rest of the server
	// configuration is irrelevant to deriving a token.
	if err := conf.BindFlags(cmd.Flags(), config.EnrolmentOptions); err != nil {
		return nil, err
	}

	return cmd, nil
}
