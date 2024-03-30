package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "v0.0.3",
	Use:     "ddns",
	Short:   "A simple DDNS server and client",
	Long: fmt.Sprintf(`A simple DDNS server and client.

ENVIRONMENT VARIABLES

%s

`, getEnvDocs("  ")),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
