// Hati
package main

import (
	"github.com/hati-sh/hati/cmd/hati/commands"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	commands.Execute()
}
