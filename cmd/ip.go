package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get the public IP of the current machine using the DDNS API",
	Run: func(cmd *cobra.Command, args []string) {
		server := "http://localhost:3345"
		if v := os.Getenv("DDNS_API_SERVER"); v != "" {
			server = v
		}

		agent := &ddns.Agent{
			ServerAddress: server,
			APIKey:        os.Getenv("DDNS_API_KEY"),
		}

		ip, err := agent.DetermineIP()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		fmt.Println(ip)
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
