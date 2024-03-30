package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update domain [ip]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Update the A record for a domain",
	Long: `Update the A record for a domain. If an IP is not provided, "auto" will
be sent in the request.

See "ddns help" for a list of supported environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := &Config{}
		if err := c.Init(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
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

		updated, err := c.Agent.UpdateIP(domain, ip)
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
