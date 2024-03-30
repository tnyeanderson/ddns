package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get the public IP of the current machine using the DDNS API",
	Run: func(cmd *cobra.Command, args []string) {
		c := &Config{}
		if err := c.Init(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		ip, err := c.Agent.DetermineIP()
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
