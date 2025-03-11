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

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/logger"
)

type SubmitFileGitHubService struct {
	logger      logger.Logger
	client      *github.Client
	callerInput functions.SubmitFilesServiceInput
}

func NewSubmitFileGitHubService(
	logger logger.Logger,
	client *github.Client,
	callerInput functions.SubmitFilesServiceInput,
) functions.SubmitFilesService {
	return SubmitFileGitHubService{
		logger:      logger,
		client:      client,
		callerInput: callerInput,
	}
}

func (s SubmitFileGitHubService) SubmitFiles(input functions.SubmitFilesInput) (submitFileOut functions.SubmitFilesOutput, _ error) {
	errorf := func(format string, a ...any) error {
		return fmt.Errorf("submit file service: "+format, a...)
	}
	var err error
	ctx := context.Background()

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
	if head.Name().Short() == s.callerInput.BaseBranch {
		return submitFileOut, errorf("cannot submit in the base branch. create and switch to a new branch")
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

	ref, err := repo.Head()
	if err != nil {
		return submitFileOut, errorf("failed to get HEAD: %w", err)
	}
	currentBranch := ref.Name().Short()
	prBranch := currentBranch

	s.logger.Debug(fmt.Sprintf("created PR parameter: name=%s, email=%s, base-branch=%s branch=%s\n",
		s.callerInput.GitName, s.callerInput.GitEmail, s.callerInput.BaseBranch, currentBranch))
	pr, _, err := s.client.PullRequests.Create(ctx, s.callerInput.GitHubOwner, s.callerInput.Repository, &github.NewPullRequest{
		Title: &input.CommitMessageShort,
		Head:  &prBranch,
		Base:  &s.callerInput.BaseBranch,
		Body:  &input.PullRequestContent,
	})
	if err != nil {
		return submitFileOut, errorf("failed to create PR: %w", err)
	}

	if len(s.callerInput.PRLabels) > 0 {
		if _, _, err = s.client.Issues.AddLabelsToIssue(
			ctx,
			s.callerInput.GitHubOwner,
			s.callerInput.Repository,
			*pr.Number,
			s.callerInput.PRLabels); err != nil {
			return submitFileOut, errorf("failed to add labels(%s) to PR: %w", s.callerInput.PRLabels, err)
		}
	}

	// checkout to the base branch
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(s.callerInput.BaseBranch),
		Keep:   false,
		Force:  true,
		Create: false,
	}); err != nil {
		return submitFileOut, fmt.Errorf("failed to checkout branch %s: %w", s.callerInput.BaseBranch, err)
	}

	return functions.SubmitFilesOutput{
		Message: fmt.Sprintf("success creating pull request.\ncreated pull request number: %d\nbranch: %s.\n switched %s branch.",
			*pr.Number, prBranch, s.callerInput.BaseBranch),
		PushedBranch:      prBranch,
		PullRequestNumber: *pr.Number,
	}, nil
}

// NopSubmitFileService implements functions.SubmitFilesService as a no-op service.
type NopSubmitFileService struct{}

// SubmitFiles is a no-op implementation of the SubmitFilesService interface.
func (s NopSubmitFileService) SubmitFiles(input functions.SubmitFilesInput) (functions.SubmitFilesOutput, error) {
	return functions.SubmitFilesOutput{
		Message:           "NopSubmitFileService: operation skipped",
		PushedBranch:      "",
		PullRequestNumber: -1,
	}, nil
}
