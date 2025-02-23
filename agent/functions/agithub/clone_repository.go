package agithub

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"

	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/logger"
)

func CloneRepository(lo logger.Logger, conf config.Config, workRepository string) error {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		lo.Error("GITHUB_TOKEN is not set")
		return fmt.Errorf("GITHUB_TOKEN is not set\n")
	}

	lo.Info("cloning repository...\n")
	if _, err := git.PlainClone(conf.WorkDir, false, &git.CloneOptions{
		URL: fmt.Sprintf("https://oauth2:%s@github.com/%s/%s.git",
			token, conf.Agent.GitHub.Owner, workRepository),
		Depth: 1,
	}); err != nil {
		return err
	}
	lo.Info("cloned repository successfully\n")

	return nil
}
