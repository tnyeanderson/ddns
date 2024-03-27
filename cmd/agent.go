package cmd

import (
	"log/slog"
	"net"
	"os"

	"github.com/spf13/cobra"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent domain [ip]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Run the DDNS update agent",
	Long: `Update the A record for a given domain. If an IP is not provided, "auto" will
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

		ip := "auto"
		if len(args) == 2 {
			ip = args[1]
			if n := net.ParseIP(ip); n == nil {
				slog.Error("ip is not valid", "ip", ip)
				os.Exit(1)
			}
		}

		if err := agent.UpdateIP(args[0], ip); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
