package cmd

import (
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent domain ip",
	Args:  cobra.MinimumNArgs(1),
	Short: "DDNS update agent",
	Long:  "Used to update the A record for a given domain",
	Run: func(cmd *cobra.Command, args []string) {
		server := "http://localhost:8989"
		if v := os.Getenv("DDNS_SERVER_ADDR"); v != "" {
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
				log.Fatalf("not a valid ip: %s", ip)
			}
		}

		if err := agent.UpdateIP(args[0], ip); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// agentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
