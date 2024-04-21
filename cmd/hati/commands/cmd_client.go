package commands

import (
	"fmt"

	"github.com/hati-sh/hati/core"
	"github.com/spf13/cobra"
)

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Connect as a client to hati server",
	Long:  `Connect as a client to hati server, provide server address eg: hati client localhost:4242`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(">>> client")
		fmt.Println(args)

		client, err := core.NewClientTcp("0.0.0.0", "4242")
		if err != nil {
			panic(err)
		}

		if err := client.Connect(); err != nil {
			panic(err)
		}
	},
}
