package main

import (
	"fmt"

	"github.com/google/go-github/v66/github"

	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/functions"
	"github.com/clover0/issue-agent/functions/agithub"
	"github.com/clover0/issue-agent/logger"
)

func main() {
	conf, err := config.Load("")
	if err != nil {
		panic(err)
	}

	lo := logger.NewPrinter(conf.LogLevel)

	functions.InitializeFunctions(
		*conf.Agent.GitHub.NoSubmit,
		agithub.NewGitHubService(conf.Agent.GitHub.Owner, "test", github.NewClient(nil).WithAuthToken(""), lo),
		conf.Agent.AllowFunctions,
	)

	out := "Functions List\n"
	for _, f := range functions.AllFunctions() {
		out += fmt.Sprintf("%s: %s\n", f.Name, f.Description)
		for propKey, values := range f.Parameters["properties"].(map[string]any) {
			propValues, ok := values.(map[string]any)
			if !ok {
				lo.Error("failed to get properties\n")
				return
			}
			out += fmt.Sprintf("    %s\n", propKey)
			out += fmt.Sprintf("        %s\n", propValues["description"])
			out += "\n"
		}
	}

	lo.Info(out)
}
