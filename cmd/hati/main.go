// Hati
package main

import (
	"runtime"

	"github.com/hati-sh/hati/cmd/hati/commands"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	runtime.GOMAXPROCS(runtime.NumCPU())

	commands.Execute()
}
