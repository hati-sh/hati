package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hati-sh/hati/core"
	"github.com/spf13/cobra"
)

const VERSION = "0.1.0-dev"

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "start hati",
	Long:  `start is for starting application.`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		tlsFlag, _ := cmd.Flags().GetString("tls")

		tlsEnabled := false

		if tlsFlag == "on" {
			tlsEnabled = true
		}

		config := &core.Config{
			ServerTcp: &core.ServerTcpConfig{
				Host:       host,
				Port:       port,
				TlsEnabled: tlsEnabled,
			},
		}

		ctx := context.Background()
		hati := core.NewHati(ctx, config)

		if err := hati.Start(); err != nil {
			panic(err)
		}

		var osSignal chan os.Signal = make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

		for {
			select {
			case <-osSignal:
				fmt.Println("stop signal...")
				hati.Stop()

				os.Exit(0)
			}
		}
	},
}
