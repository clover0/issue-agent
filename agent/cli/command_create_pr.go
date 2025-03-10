package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v69/github"

	"github.com/clover0/issue-agent/agithub"
	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/core"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/models"
)

const CreatePrCommand = "create-pr"

func CreatePR(flags []string) error {
	cliIn, err := ParseCreatePRInput(flags)
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	conf, err := config.LoadDefault(isPassedConfig(cliIn.Common.Config))
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	conf = cliIn.MergeConfig(conf)

	if err := config.ValidateConfig(conf); err != nil {
		return err
	}

	lo := logger.NewPrinter(conf.LogLevel)

	if *conf.Agent.GitHub.CloneRepository {
		if err := agithub.CloneRepository(lo, conf, cliIn.WorkRepository, cliIn.BaseBranch); err != nil {
			lo.Error("failed to clone repository")
			return err
		}
	}

	// TODO: no dependency with changing directory
	if err := os.Chdir(conf.WorkDir); err != nil {
		lo.Error("failed to change directory: %s\n", err)
		return err
	}

	gh, err := newGitHub()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	ctx := context.Background()

	return core.OrchestrateAgentsByIssue(ctx, lo, conf, cliIn.BaseBranch, cliIn.WorkRepository, gh, cliIn.GithubIssueNumber, models.SelectForwarder)
}

func isPassedConfig(configPath string) bool {
	return configPath != ""
}

func newGitHub() (*github.Client, error) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return nil, fmt.Errorf("GITHUB_TOKEN is not set")
	}
	return github.NewClient(nil).WithAuthToken(token), nil
}
