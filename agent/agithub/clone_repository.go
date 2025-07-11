package agithub

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/clover0/issue-agent/logger"
)

func CloneRepository(lo logger.Logger, owner string, workRepository string, refBranch string) error {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		lo.Error("GITHUB_TOKEN is not set")
		return fmt.Errorf("GITHUB_TOKEN is not set")
	}

	lo.Info("cloning repository...\n")
	if _, err := git.PlainClone(path.Join(".", workRepository), false, &git.CloneOptions{
		URL: fmt.Sprintf("https://oauth2:%s@github.com/%s/%s.git",
			token, owner, workRepository),
		Depth:         1,
		ReferenceName: plumbing.ReferenceName(refBranch),
	}); err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return fmt.Errorf("branch %s not found in repository %s/%s", refBranch, owner, workRepository)
		}
		return err
	}
	lo.Info("cloned repository successfully\n")

	return nil
}
