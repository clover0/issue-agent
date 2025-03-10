package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/clover0/issue-agent/logger"
)

const HelpCommand = "help"

func Help(lo logger.Logger) {
	msg := `Usage
  issue-agent <command> [flags]
Command and Flags  
  help: Show usage of commands and flags
  version: Show version of issue-agent CLI
`
	createPRFlags, _ := CreatePRFlags()
	reactFlags, _ := ReactFlags()

	msg += fmt.Sprintf("  %s:\n", CreatePrCommand)
	msg += "    Usage:\n"
	msg += fmt.Sprintf("      %s GITHUB_OWNER/REPOSITORY/issues/NUMBER [flags]\n", CreatePrCommand)
	msg += "    Flags:\n"

	createPRFlags.VisitAll(func(flg *flag.Flag) {
		msg += fmt.Sprintf("    --%s\n", flg.Name)
		msg += IndentMultiLine(flg.Usage, "      ")
		msg += "\n"
	})

	msg += fmt.Sprintf("  %s:\n", ReactCommand)
	msg += "    Usage:\n"
	msg += fmt.Sprintf("      %s OWNER/REPO/issues/comments/COMMENT_ID [flags]\n", ReactCommand)
	msg += "    Flags:\n"
	reactFlags.VisitAll(func(flg *flag.Flag) {
		msg += fmt.Sprintf("    --%s\n", flg.Name)
		msg += IndentMultiLine(flg.Usage, "      ")
		msg += "\n"
	})

	lo.Info(msg)
}

func IndentMultiLine(str string, indent string) string {
	lines := strings.Split(str, "\n")
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = indent + line
	}

	return strings.Join(out, "\n")
}
