package main

import (
	"os"

	"github.com/setcreed/store-core/cmd/app"
)

func main() {
	cmd := app.NewServerCommand("version")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
