package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "bob",
	Short:         "Bob builds images",
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute runs the top level command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
