package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE:    sync,
	}
	return cmd
}

func sync(cmd *cobra.Command, args []string) error {
	// ctx := cmd.Context()

	// convert args
	// reader, err := workload.NewReader(args)
	// if err != nil {
	// 	return err
	// }

	// new sources
	// sources := reader.Sources

	// new destinations
	// destinations := reader.Destinations

	// new transformations
	// transformers := reader.Transformers

	// start sync
	// for

	return nil
}
