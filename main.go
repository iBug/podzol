package main

import (
	"os"

	"github.com/ustclug/podzol/cmd"
)

func main() {
	if cmd.Execute() != nil {
		os.Exit(1)
	}
}
