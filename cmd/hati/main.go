// Hati
package main

import (
	"fmt"
	"runtime"

	"github.com/hati-sh/hati/cmd/hati/commands"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	fmt.Println(runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	commands.Execute()
}
