package cli

import (
	"fmt"
	"os"

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
	case VersionCommand:
		return Version()
	case CreatePrCommand:
		return CreatePR(others)
	case HelpCommand:
		Help(lo)
		return nil
	default:
		Help(lo)
		return fmt.Errorf("unknown command: %s", command)
	}
}
