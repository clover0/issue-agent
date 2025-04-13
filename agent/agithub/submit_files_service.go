package agithub

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v70/github"

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
	ctx := context.Background()

	if err := s.validateGitConfig(); err != nil {
		return submitFileOut, errorf("failed to validate git config: %w", err)
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		return submitFileOut, errorf("failed to open repository: %w", err)
	}

	if err := s.guardPushToBaseBranch(repo); err != nil {
		return submitFileOut, errorf("failed on guard push to base branch: %w", err)
	}

	if err := s.setGitConfig(repo); err != nil {
		return submitFileOut, errorf("failed to set git config: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return submitFileOut, errorf("failed to get worktree: %w", err)
	}

	if _, err := wt.Add("./"); err != nil {
		return submitFileOut, errorf("failed to add files: %w", err)
	}

	if err := s.resetSymlink(wt); err != nil {
		return submitFileOut, errorf("failed to reset symlink: %w", err)
	}

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
	pushedBranch := ref.Name().Short()

	s.logger.Debug(fmt.Sprintf("created PR parameter: name=%s, email=%s, base-branch=%s branch=%s\n",
		s.callerInput.GitName, s.callerInput.GitEmail, s.callerInput.BaseBranch, pushedBranch))
	pr, _, err := s.client.PullRequests.Create(ctx, s.callerInput.GitHubOwner, s.callerInput.Repository, &github.NewPullRequest{
		Title: &input.CommitMessageShort,
		Head:  &pushedBranch,
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

	if len(s.callerInput.Reviewers) > 0 || len(s.callerInput.TeamReviewers) > 0 {
		if err := s.reviewRequest(ctx, pr, s.callerInput.Reviewers, s.callerInput.TeamReviewers); err != nil {
			return submitFileOut, errorf("failed to request reviewers: %w", err)
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
			*pr.Number, pushedBranch, s.callerInput.BaseBranch),
		PushedBranch:      pushedBranch,
		PullRequestNumber: *pr.Number,
	}, nil
}

// validateGitConfig validate git config.
func (s SubmitFileGitHubService) validateGitConfig() error {
	if s.callerInput.GitEmail == "" {
		return fmt.Errorf("git email is not set")
	}
	if s.callerInput.GitName == "" {
		return fmt.Errorf("git name is not set")
	}

	return nil
}

// guardPushToBaseBranch guard pushing to base branch.
// If the current branch is the base branch, return an error.
func (s SubmitFileGitHubService) guardPushToBaseBranch(repo *git.Repository) error {
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}
	if head.Name().Short() == s.callerInput.BaseBranch {
		return fmt.Errorf("cannot submit in the base branch. create and switch to a new branch")
	}

	return nil
}

// setGitConfig set git config.
func (s SubmitFileGitHubService) setGitConfig(repo *git.Repository) error {
	cfg, err := repo.Config()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	cfg.User.Email = s.callerInput.GitEmail
	cfg.User.Name = s.callerInput.GitName

	if err := repo.SetConfig(cfg); err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	return nil
}

// resetSymlink reset symlink from git staging.
// Because go-git's file system behavior causes symlinks to be relative paths, resulting in extra diffs.
func (s SubmitFileGitHubService) resetSymlink(wt *git.Worktree) error {
	statuses, err := wt.Status()
	if err != nil {
		return fmt.Errorf("failed to get worktree status: %w", err)
	}

	for path, status := range statuses {
		if status.Staging != git.Modified {
			continue
		}
		f, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		if f.Mode()&os.ModeSymlink != 0 {
			s.logger.Debug(fmt.Sprintf("reset symlink: %s\n", path))
			if err := wt.Reset(&git.ResetOptions{Files: []string{path}}); err != nil {
				return fmt.Errorf("failed to reset symlink: %w", err)
			}
		}
	}
	s.logger.Info(statuses.String())

	return nil
}

func (s SubmitFileGitHubService) reviewRequest(ctx context.Context, pr *github.PullRequest, reviewers []string, teamReviewers []string) error {
	if _, resp, err := s.client.PullRequests.RequestReviewers(
		ctx,
		s.callerInput.GitHubOwner,
		s.callerInput.Repository,
		*pr.Number,
		github.ReviewersRequest{
			Reviewers:     reviewers,
			TeamReviewers: teamReviewers,
		}); err != nil {
		if resp == nil {
			return fmt.Errorf("failed to request reviewers=%s, team_reviewers=%s to PR: %w", s.callerInput.Reviewers, s.callerInput.TeamReviewers, err)
		}

		// if the client error caused, print the response body and continue
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}
			defer resp.Body.Close()

			s.logger.Error("client error: %s\n", body)
		}
	}

	return nil
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
