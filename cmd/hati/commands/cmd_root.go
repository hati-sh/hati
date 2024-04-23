package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "hati", Short: "hati cli", Long: "", Version: VERSION}

func init() {
	rootCmd.PersistentFlags().StringP("author", "a", "", "Maciej Lisowski")
}

func Execute() {
	rootCmd.AddCommand(cmdStart, cmdClient)

	cmdStart.PersistentFlags().String("host", "0.0.0.0", "bind address for TCP server")
	cmdStart.PersistentFlags().String("port", "4242", "bind port for TCP server")

	cmdClient.PersistentFlags().String("host", "0.0.0.0", "address to connect to")
	cmdClient.PersistentFlags().String("port", "4242", "target port")

	var tlsFlagValue string

	cmdStart.Flags().StringVar(&tlsFlagValue, "tls", "off", "tls on/off, defaut: off")
	cmdClient.Flags().StringVar(&tlsFlagValue, "tls", "off", "tls on/off, defaut: off")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
