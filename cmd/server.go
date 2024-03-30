package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a DDNS server",
	Long:  "Start an HTTP server and a DNS server.",
	Run: func(cmd *cobra.Command, args []string) {
		c := &Config{}
		if err := c.Init(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		// Load domains from hosts file
		if err := c.Server.Load(); err != nil {
			// Log, but don't exit here. Failing to load from the hosts file is not
			// fatal, and should not stop the DNS server from starting. It's better
			// to start up and wait for the next update than to crash and have no
			// chance of fulfilling requests at all.
			slog.Error(err.Error())
		}

		// Start the server
		if err := c.Server.Listen(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
