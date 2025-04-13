package main

import (
	"fmt"

	"github.com/google/go-github/v70/github"

	"github.com/clover0/issue-agent/agithub"
	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/logger"
)

func main() {
	conf, err := config.Load("")
	if err != nil {
		panic(err)
	}

	lo := logger.NewPrinter(conf.LogLevel)

	ghClient := github.NewClient(nil).WithAuthToken("")
	submitService := agithub.NopSubmitFileService{}
	revisionService := agithub.NopSubmitRevisionService{}

	functions.InitializeFunctions(
		*conf.Agent.GitHub.NoSubmit,
		agithub.NewGitHubService(conf.Agent.GitHub.Owner, "repo", ghClient, lo),
		submitService,
		revisionService,
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
