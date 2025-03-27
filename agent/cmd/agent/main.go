package main

import (
	"os"

	"github.com/clover0/issue-agent/cli/command"
	"github.com/clover0/issue-agent/logger"
)

func main() {
	// TODO:
	//lo := logger.NewDefaultLogger()
	lo := logger.NewPrinter("info")
	lo.Info("start agent in container...\n")

	if err := command.Execute(); err != nil {
		lo.Error("failed to execute command: %s\n", err)
		os.Exit(1)
	}
}
