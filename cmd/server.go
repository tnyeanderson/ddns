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
		c := Config{}
		if err := c.Init(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		s := c.Server

		//// Overwrite config values with env vars, if set
		//if v := os.Getenv(EnvServerHTTPListener); v != "" {
		//	s.HTTPListener = v
		//}

		//if v := os.Getenv(EnvServerDNSListener); v != "" {
		//	s.DNSListener = v
		//}

		//if key := os.Getenv(EnvServerAPIKey); key != "" {
		//	var r *regexp.Regexp
		//	if pattern := os.Getenv(EnvServerAPIKeyRegex); pattern != "" {
		//		r = regexp.MustCompile(pattern)
		//	}
		//	s.Allow(key, r)
		//}

		//if v := os.Getenv(EnvServerHostsFile); v != "" {
		//	s.HostsFile = v
		//}

		// Load domains from hosts file
		if err := s.Load(); err != nil {
			// Log, but don't exit here. Failing to load from the hosts file is not
			// fatal, and should not stop the DNS server from starting. It's better
			// to start up and wait for the next update than to crash and have no
			// chance of fulfilling requests at all.
			slog.Error(err.Error())
		}

		// Start the server
		if err := s.Listen(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
