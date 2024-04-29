package commands

import (
	"fmt"
	"github.com/hati-sh/hati/common"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "hati", Short: "hati cli", Long: "", Version: common.VERSION}

func init() {
	rootCmd.PersistentFlags().StringP("author", "a", "", "Maciej Lisowski")
}

func Execute() {
	rootCmd.AddCommand(cmdStart, cmdClient)

	cmdStart.PersistentFlags().String("tcp-host", common.DEFAULT_TCP_BIND_HOST, "bind address for TCP server")
	cmdStart.PersistentFlags().String("tcp-port", common.DEFAULT_TCP_BIND_PORT, "bind port for TCP server")
	cmdStart.PersistentFlags().String("rpc-host", common.DEFAULT_RPC_BIND_HOST, "bind address for JSON-RPC server")
	cmdStart.PersistentFlags().String("rpc-port", common.DEFAULT_RPC_BIND_PORT, "bind port for JSON-RPC server")
	cmdStart.PersistentFlags().String("data-dir", "", "absolute path to directory where store data")

	cmdClient.PersistentFlags().String("tcp-host", common.DEFAULT_TCP_BIND_HOST, "address to connect to")
	cmdClient.PersistentFlags().String("tcp-port", common.DEFAULT_TCP_BIND_PORT, "target port")

	var tlsFlagValue bool
	var tcpFlagValue bool
	var rpcFlagValue bool

	cmdStart.Flags().BoolVar(&tcpFlagValue, "tcp", false, "tcp server true/false, default: false")
	cmdStart.Flags().BoolVar(&tlsFlagValue, "tcp-tls", false, "tls for tcp server on/off, default: false")
	cmdStart.Flags().BoolVar(&rpcFlagValue, "rpc", false, "rpc server true/false, default: false")

	cmdClient.Flags().BoolVar(&tlsFlagValue, "tcp-tls", false, "tls on/off, default: false")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
