package cmd

import (
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a DDNS server",
	Long:  "This will start an HTTP server and a DNS server",
	Run: func(cmd *cobra.Command, args []string) {
		s := &ddns.Server{}

		if v := os.Getenv("DDNS_SERVER_HTTP_LISTENER"); v != "" {
			s.HTTPListener = v
		}

		if v := os.Getenv("DDNS_SERVER_DNS_LISTENER"); v != "" {
			s.DNSListener = v
		}

		s.Allow("mytesttoken", regexp.MustCompile("hello.world"))
		if err := s.Listen(); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
