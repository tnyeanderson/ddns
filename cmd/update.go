package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/spf13/cobra"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

var updateCmd = &cobra.Command{
	Use:   "update domain [ip]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Update the A record for a domain",
	Long: `Update the A record for a domain. If an IP is not provided, "auto" will
be sent in the request.`,
	Run: func(cmd *cobra.Command, args []string) {
		server := "http://localhost:3345"
		if v := os.Getenv("DDNS_API_SERVER"); v != "" {
			server = v
		}

		agent := &ddns.Agent{
			ServerAddress: server,
			APIKey:        os.Getenv("DDNS_API_KEY"),
		}

		domain := args[0]

		ip := "auto"
		if len(args) == 2 {
			ip = args[1]
			if n := net.ParseIP(ip); n == nil {
				slog.Error("ip is not valid", "ip", ip)
				os.Exit(1)
			}
		}

		updated, err := agent.UpdateIP(domain, ip)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		if updated {
			slog.Info(fmt.Sprintf("updated dns entry for %s to %s", domain, ip))
		} else {
			slog.Info(fmt.Sprintf("dns entry already correct for %s: %s", domain, ip))
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
