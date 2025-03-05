package agithub

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v69/github"

	"github.com/clover0/issue-agent/functions"
	"github.com/clover0/issue-agent/logger"
)

// TODO: move to GitHub service
type SubmitFileGitHubService struct {
	owner      string
	repository string
	client     *github.Client
	logger     logger.Logger
}

func NewSubmitFileGitHubService(
	owner string,
	repository string,
	client *github.Client,
	logger logger.Logger,
) functions.SubmitFilesService {
	return SubmitFileGitHubService{
		owner:      owner,
		repository: repository,
		client:     client,
		logger:     logger,
	}
}

// TODO: move to GitHub service
func (s SubmitFileGitHubService) Caller(
	ctx context.Context,
	callerInput functions.SubmitFilesServiceInput,
) functions.SubmitFilesCallerType {
	errorf := func(format string, a ...any) error {
		return fmt.Errorf("submit file service: "+format, a...)
	}

	return func(input functions.SubmitFilesInput) (submitFileOut functions.SubmitFilesOutput, _ error) {
		var err error

		// TODO: validation before this caller
		if callerInput.GitEmail == "" {
			return submitFileOut, errorf("git email is not set")
		}
		if callerInput.GitName == "" {
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
		if head.Name().Short() == callerInput.BaseBranch {
			return submitFileOut, errorf("cannot submit in the base branch. create and switch to a new branch")
		}

		cfg, err := repo.Config()
		if err != nil {
			return submitFileOut, err
		}

		cfg.User.Email = callerInput.GitEmail
		cfg.User.Name = callerInput.GitName

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

		status, err := wt.Status()
		if err != nil {
			return submitFileOut, errorf("failed to get worktree status: %w", err)
		}

		// reset symlink becauseb go-git's file system behavior causes symlinks to be relative paths, resulting in extra diffs.
		for path := range status {
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
		s.logger.Info(status.String())

		if _, err := wt.Commit(
			fmt.Sprintf("%s\n\n%s", input.CommitMessageShort, input.CommitMessageDetail),
			&git.CommitOptions{
				Author: &object.Signature{
					Name:  callerInput.GitName,
					Email: callerInput.GitEmail,
					When:  time.Now(),
				},
			}); err != nil {
			return submitFileOut, errorf("failed to commit: %w", err)
		}

		if repo.Push(&git.PushOptions{RemoteName: "origin"}) != nil {
			return submitFileOut, errorf("failed to push: %w", err)
		}

		ref, err := repo.Head()
		if err != nil {
			return submitFileOut, errorf("failed to get HEAD: %w", err)
		}
		currentBranch := ref.Name().Short()
		prBranch := currentBranch

		s.logger.Debug(fmt.Sprintf("created PR parameter: name=%s, email=%s, base-branch=%s branch=%s\n",
			callerInput.GitName, callerInput.GitEmail, callerInput.BaseBranch, currentBranch))
		pr, _, err := s.client.PullRequests.Create(ctx, s.owner, s.repository, &github.NewPullRequest{
			Title: &input.CommitMessageShort,
			Head:  &prBranch,
			Base:  &callerInput.BaseBranch,
			Body:  &input.PullRequestContent,
		})
		if err != nil {
			return submitFileOut, errorf("failed to create PR: %w", err)
		}

		if len(callerInput.PRLabels) > 0 {
			if _, _, err = s.client.Issues.AddLabelsToIssue(
				ctx,
				s.owner,
				s.repository,
				*pr.Number,
				callerInput.PRLabels); err != nil {
				return submitFileOut, errorf("failed to add labels(%s) to PR: %w", callerInput.PRLabels, err)
			}
		}

		// checkout to the base branch
		if err := wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(callerInput.BaseBranch),
			Keep:   false,
			Force:  true,
			Create: false,
		}); err != nil {
			return submitFileOut, fmt.Errorf("failed to checkout branch %s: %w", callerInput.BaseBranch, err)
		}

		return functions.SubmitFilesOutput{
			Message: fmt.Sprintf("success creating pull request.\ncreated pull request number: %d\nbranch: %s.\n switched %s branch.",
				*pr.Number, prBranch, callerInput.BaseBranch),
			PushedBranch:      prBranch,
			PullRequestNumber: *pr.Number,
		}, nil
	}
}
