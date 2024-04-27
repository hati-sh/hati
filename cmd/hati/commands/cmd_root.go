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
	cmdStart.PersistentFlags().String("rpc-host", "0.0.0.0", "bind address for JSON-RPC server")
	cmdStart.PersistentFlags().String("rpc-port", "6767", "bind port for JSON-RPC server")

	cmdClient.PersistentFlags().String("host", "0.0.0.0", "address to connect to")
	cmdClient.PersistentFlags().String("port", "4242", "target port")

	var tlsFlagValue bool
	var rpcFlagValue bool

	cmdStart.Flags().BoolVar(&tlsFlagValue, "tls", false, "tls on/off, defaut: off")
	cmdStart.Flags().BoolVar(&rpcFlagValue, "rpc", false, "rpc server true/false, defaut: false")

	cmdClient.Flags().BoolVar(&tlsFlagValue, "tls", false, "tls on/off, defaut: off")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
