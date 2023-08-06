package main

import "github.com/ustclug/podzol/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
