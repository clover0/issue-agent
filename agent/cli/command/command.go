package command

import (
	"fmt"
	"os"

	"github.com/clover0/issue-agent/cli/command/createpr"
	"github.com/clover0/issue-agent/cli/command/help"
	"github.com/clover0/issue-agent/cli/command/react"
	"github.com/clover0/issue-agent/cli/command/version"
	"github.com/clover0/issue-agent/logger"
)

const noCommand = "no-command"

// Parse parses the command and others from os.Args
// issue-agent <command> others
func Parse() (command string, others []string) {
	if len(os.Args) < 2 {
		return noCommand, []string{}
	}

	return os.Args[1], os.Args[2:]
}

func Execute() error {
	command, others := Parse()

	lo := logger.NewPrinter("info")
	switch command {
	case version.VersionCommand:
		return version.Version()
	case createpr.CreatePrCommand:
		return createpr.CreatePR(others)
	case react.ReactCommand:
		return react.React(others)
	case help.HelpCommand:
		help.Help(lo)
		return nil
	default:
		help.Help(lo)
		return fmt.Errorf("unknown command: %s", command)
	}
}
