package commands

import (
	"github.com/hati-sh/hati/core"
	"github.com/spf13/cobra"
)

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Connect to hati server",
	Long:  `Connect to hati server, provide server address eg: localhost:4242`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		tlsFlag, _ := cmd.Flags().GetString("tls")

		tlsEnabled := false
		if tlsFlag == "on" {
			tlsEnabled = true
		}

		client, err := core.NewClientTcp(host, port, tlsEnabled)
		if err != nil {
			panic(err)
		}

		if err := client.Connect(); err != nil {
			panic(err)
		}
	},
}
