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
	rootCmd.AddCommand(cmdStart)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
