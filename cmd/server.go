package cmd

import (
	"log/slog"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

// serverCmd represents the server command
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

		if token := os.Getenv("DDNS_SERVER_TOKEN"); token != "" {
			var r *regexp.Regexp
			if pattern := os.Getenv("DDNS_SERVER_TOKEN_REGEX"); pattern != "" {
				r = regexp.MustCompile(pattern)
			}
			s.Allow(token, r)
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
