package agithub

import (
	"fmt"
	"os"
	"os/exec"

	"github/clover0/github-issue-agent/config/cli"
	"github/clover0/github-issue-agent/logger"
)

func CloneRepository(lo logger.Logger, cliIn cli.IssueInputs) error {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		lo.Error("GITHUB_TOKEN is not set")
		return fmt.Errorf("GITHUB_TOKEN is not set")
	}
	lo.Info("cloning repository...\n")
	cmd := exec.Command("git", "clone", "--depth", "1",
		fmt.Sprintf("https://oauth2:%s@github.com/%s/%s.git", token, cliIn.RepositoryOwner, cliIn.Repository),
		cliIn.AgentWorkDir,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		lo.Error(string(output))
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	lo.Info("cloned repository successfully")
	return nil
}
