package cmd

import (
	"log/slog"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a DDNS server",
	Long:  "Start an HTTP server and a DNS server.",
	Run: func(cmd *cobra.Command, args []string) {
		s := &ddns.Server{}

		if v := os.Getenv("DDNS_SERVER_HTTP_LISTENER"); v != "" {
			s.HTTPListener = v
		}

		if v := os.Getenv("DDNS_SERVER_DNS_LISTENER"); v != "" {
			s.DNSListener = v
		}

		if key := os.Getenv("DDNS_SERVER_API_KEY"); key != "" {
			var r *regexp.Regexp
			if pattern := os.Getenv("DDNS_SERVER_API_KEY_REGEX"); pattern != "" {
				r = regexp.MustCompile(pattern)
			}
			s.Allow(key, r)
		}

		if v := os.Getenv("DDNS_SERVER_CACHE_FILE"); v != "" {
			s.CacheFile = v
		}

		if err := s.Load(); err != nil {
			// Log, but don't exit here. Failing to load from the cache is not fatal
			// and should not stop the DNS server from starting.
			slog.Error(err.Error())
		}

		if err := s.Listen(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
