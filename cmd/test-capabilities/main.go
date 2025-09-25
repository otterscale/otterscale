package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	oscmd "github.com/otterscale/otterscale/internal/cmd"
)

func main() {
	cmd := &cobra.Command{
		Use:   "test-capabilities",
		Short: "Test capabilities command",
	}
	
	cmd.AddCommand(oscmd.NewCapabilities())
	
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}