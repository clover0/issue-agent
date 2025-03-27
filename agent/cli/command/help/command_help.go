package help

import (
	"flag"
	"fmt"
	"strings"

	"github.com/clover0/issue-agent/cli/command/createpr"
	"github.com/clover0/issue-agent/cli/command/react"
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
	createPRFlags, _ := createpr.CreatePRFlags()
	reactFlags, _ := react.ReactFlags()

	msg += fmt.Sprintf("  %s:\n", createpr.CreatePrCommand)
	msg += "    Usage:\n"
	msg += fmt.Sprintf("      %s GITHUB_OWNER/REPOSITORY/issues/NUMBER [flags]\n", createpr.CreatePrCommand)
	msg += "    Flags:\n"

	createPRFlags.VisitAll(func(flg *flag.Flag) {
		msg += fmt.Sprintf("    --%s\n", flg.Name)
		msg += IndentMultiLine(flg.Usage, "      ")
		msg += "\n"
	})

	msg += fmt.Sprintf("  %s:\n", react.ReactCommand)
	msg += "    Usage:\n"
	msg += fmt.Sprintf("      %s RESOURCE_FORMAT [flags]\n", react.ReactCommand)
	msg += "    RESOURCE_FORMAT:\n"
	msg += "        issue_comment(pull request comment): OWNER/REPO/issues/comments/COMMENT_ID\n"
	msg += "        pull_request_review_comment: OWNER/REPO/pulls/comments/COMMENT_ID\n"
	msg += "    Example:\n"
	msg += "       react owner/example/issues/comments/123456 [flags]\n"
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
