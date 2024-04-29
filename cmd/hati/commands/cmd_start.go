package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/hati-sh/hati/common/logger"
	"github.com/hati-sh/hati/core"
	"github.com/spf13/cobra"
)

const VERSION = "0.1.0-dev"

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "start hati",
	Long:  `start is for starting application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log("++ hati ++ v" + VERSION)

		cpuNum, _ := cmd.Flags().GetInt("cpu-num")

		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")

		tcpFlag, _ := cmd.Flags().GetBool("tcp")
		tlsFlag, _ := cmd.Flags().GetBool("tls")
		rpcFlag, _ := cmd.Flags().GetBool("rpc")

		rpcHost, _ := cmd.Flags().GetString("rpc-host")
		rpcPort, _ := cmd.Flags().GetString("rpc-port")

		dataDir, _ := cmd.Flags().GetString("data-dir")

		tcpEnabled := false
		tlsEnabled := false
		rpcEnabled := false

		if cpuNum == 0 {
			cpuNum = runtime.NumCPU()
		}

		runtime.GOMAXPROCS(cpuNum)
		logger.Debug("Max CPU num: " + fmt.Sprint(cpuNum))

		if host == "" {
			host = "0.0.0.0"
		}

		if port == "" {
			port = "4242"
		}

		if tcpFlag {
			tcpEnabled = true
		}

		if tlsFlag {
			tlsEnabled = true
		}

		if rpcFlag {
			rpcEnabled = true
		}

		config := &core.Config{
			ServerTcp: &core.TcpServerConfig{
				Host:       host,
				Port:       port,
				Enabled:    tcpEnabled,
				TlsEnabled: tlsEnabled,
			},
			ServerRpc: &core.RpcServerConfig{
				Host:    rpcHost,
				Port:    rpcPort,
				Enabled: rpcEnabled,
			},
		}

		fmt.Printf("dataDir: %s\n\n", dataDir)

		ctx := context.Background()
		hati := core.NewHati(ctx, config)

		if err := hati.Start(); err != nil {
			panic(err)
		}

		var osSignal = make(chan os.Signal, 1)
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
