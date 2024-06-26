package commands

import (
	"context"
	"fmt"
	"github.com/hati-sh/hati/common"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/hati-sh/hati/common/logger"
	"github.com/hati-sh/hati/core"
	"github.com/spf13/cobra"
)

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "start hati",
	Long:  `start is for starting application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log("++ hati ++ v" + common.VERSION)

		cpuNum, _ := cmd.Flags().GetInt("cpu-num")

		host, _ := cmd.Flags().GetString("tcp-host")
		port, _ := cmd.Flags().GetString("tcp-port")

		tcpFlag, _ := cmd.Flags().GetBool("tcp")
		tlsFlag, _ := cmd.Flags().GetBool("tcp-tls")
		rpcFlag, _ := cmd.Flags().GetBool("rpc")

		rpcHost, _ := cmd.Flags().GetString("rpc-host")
		rpcPort, _ := cmd.Flags().GetString("rpc-port")

		dataDir, _ := cmd.Flags().GetString("data-dir")

		tcpEnabled := false
		tcpTlsEnabled := false
		rpcEnabled := false

		if cpuNum == 0 {
			cpuNum = runtime.NumCPU()
		}

		runtime.GOMAXPROCS(cpuNum)
		logger.Debug("Max CPU num: " + fmt.Sprint(cpuNum))

		if tcpFlag {
			tcpEnabled = true
		}

		if tlsFlag {
			tcpTlsEnabled = true
		}

		if rpcFlag {
			rpcEnabled = true
		}

		config := &core.Config{
			ServerTcp: &core.TcpServerConfig{
				Host:       host,
				Port:       port,
				Enabled:    tcpEnabled,
				TlsEnabled: tcpTlsEnabled,
			},
			ServerRpc: &core.RpcServerConfig{
				Host:    rpcHost,
				Port:    rpcPort,
				Enabled: rpcEnabled,
			},
		}

		if dataDir == "" {
			dataDir = common.DEFAULT_DATA_DIR
		}

		if _, err := os.Stat(dataDir); err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(dataDir, 0744)
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
		config.DataDir = dataDir

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
				logger.Log("shutting down, wait...")

				hati.Stop()

				os.Exit(0)
			}
		}
	},
}
