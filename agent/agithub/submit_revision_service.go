package agithub

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v71/github"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/logger"
)

type SubmitRevisionGitHubService struct {
	logger      logger.Logger
	client      *github.Client
	callerInput functions.SubmitRevisionServiceInput
}

func NewSubmitRevisionGitHubService(
	logger logger.Logger,
	client *github.Client,
	callerInput functions.SubmitRevisionServiceInput,
) functions.SubmitRevisionService {
	return SubmitRevisionGitHubService{
		logger:      logger,
		client:      client,
		callerInput: callerInput,
	}
}

func (s SubmitRevisionGitHubService) SubmitRevision(input functions.SubmitRevisionInput) (submitFileOut functions.SubmitRevisionOutput, _ error) {
	errorf := func(format string, a ...any) error {
		return fmt.Errorf("submit revision service: "+format, a...)
	}
	var err error

	// TODO: validation before this caller
	if s.callerInput.GitEmail == "" {
		return submitFileOut, errorf("git email is not set")
	}
	if s.callerInput.GitName == "" {
		return submitFileOut, errorf("git  name is not set")
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		return submitFileOut, errorf("failed to open repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return submitFileOut, errorf("failed to get HEAD: %w", err)
	}
	if head.Name().Short() != s.callerInput.WorkBranch {
		return submitFileOut, errorf("current branch is not work branch: %s", head.Name().Short())
	}

	cfg, err := repo.Config()
	if err != nil {
		return submitFileOut, err
	}

	cfg.User.Email = s.callerInput.GitEmail
	cfg.User.Name = s.callerInput.GitName

	if err := repo.SetConfig(cfg); err != nil {
		return submitFileOut, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return submitFileOut, errorf("failed to get worktree: %w", err)
	}

	if _, err := wt.Add("./"); err != nil {
		return submitFileOut, errorf("failed to add files: %w", err)
	}

	statuses, err := wt.Status()
	if err != nil {
		return submitFileOut, errorf("failed to get worktree status: %w", err)
	}

	// reset symlink because go-git's file system behavior causes symlinks to be relative paths, resulting in extra diffs.
	for path, status := range statuses {
		if status.Staging != git.Modified {
			continue
		}
		f, err := os.Lstat(path)
		if err != nil {
			return submitFileOut, fmt.Errorf("failed to open file %s: %w", path, err)
		}
		if f.Mode()&os.ModeSymlink != 0 {
			s.logger.Debug(fmt.Sprintf("reset symlink: %s\n", path))
			if err := wt.Reset(&git.ResetOptions{Files: []string{path}}); err != nil {
				return submitFileOut, errorf("failed to reset symlink: %w", err)
			}
		}
	}
	s.logger.Info(statuses.String())

	if _, err := wt.Commit(
		fmt.Sprintf("%s\n\n%s", input.CommitMessageShort, input.CommitMessageDetail),
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  s.callerInput.GitName,
				Email: s.callerInput.GitEmail,
				When:  time.Now(),
			},
		}); err != nil {
		return submitFileOut, errorf("failed to commit: %w", err)
	}

	if err := repo.Push(&git.PushOptions{RemoteName: "origin"}); err != nil {
		return submitFileOut, errorf("failed to push: %w", err)
	}

	return functions.SubmitRevisionOutput{}, nil
}

// NopSubmitRevisionService implements functions.SubmitRevisionsService as a no-op service.
type NopSubmitRevisionService struct{}

// SubmitRevision is a no-op implementation of the SubmitRevisionsService interface.
func (s NopSubmitRevisionService) SubmitRevision(_ functions.SubmitRevisionInput) (functions.SubmitRevisionOutput, error) {
	return functions.SubmitRevisionOutput{}, nil
}
