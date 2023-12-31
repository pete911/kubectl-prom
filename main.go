package main

import (
	"github.com/pete911/kubectl-prom/cmd"
	"os"
)

var Version = "dev"

func main() {
	cmd.Version = Version
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
