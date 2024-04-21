package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hati-sh/hati/core"
	"github.com/spf13/cobra"
)

const VERSION = "0.1.0-dev"

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "start hati",
	Long:  `start is for starting application.`,
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		hati := core.NewHati()
		if err := hati.Start(); err != nil {
			panic(err)
		}

		msg := core.NewMessage()
		msg.SetPayload([]byte("hello"))

		fmt.Println(msg.Bytes())
		fmt.Println(len(msg.Bytes()))
		fmt.Println(string(msg.Bytes()))

		fmt.Println("hati " + VERSION)

		var osSignal chan os.Signal = make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

		go func() {
			timer := time.NewTicker(5 * time.Second)

			for {
				select {
				case <-timer.C:
					// fmt.Println("jest")
				}
			}
		}()

		for {
			select {
			case <-osSignal:
				fmt.Println("shutting down")

				os.Exit(0)
			}
		}
	},
}
