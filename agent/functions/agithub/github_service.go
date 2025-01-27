package agithub

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v68/github"

	"github.com/clover0/issue-agent/functions"
	"github.com/clover0/issue-agent/logger"
)

type GitHubService struct {
	owner      string
	repository string
	client     *github.Client
	logger     logger.Logger
}

func NewGitHubService(
	owner string,
	repository string,
	client *github.Client,
	logger logger.Logger,
) GitHubService {
	return GitHubService{
		owner:      owner,
		repository: repository,
		client:     client,
		logger:     logger,
	}
}

func (s GitHubService) GetPullRequest(prNumber string) (functions.GetPullRequestOutput, error) {
	number, err := strconv.Atoi(prNumber)
	if err != nil {
		return functions.GetPullRequestOutput{}, fmt.Errorf("failed to convert pull request number to int: %w", err)
	}

	c := context.Background()
	pr, _, err := s.client.PullRequests.Get(c, s.owner, s.repository, number)
	if err != nil {
		return functions.GetPullRequestOutput{}, fmt.Errorf("failed to get pull request: %w", err)
	}
	diff, _, err := s.client.PullRequests.GetRaw(c, s.owner, s.repository, number, github.RawOptions{Type: github.Diff})
	if err != nil {
		return functions.GetPullRequestOutput{}, fmt.Errorf("failed to get pull request diff: %w", err)
	}

	return functions.GetPullRequestOutput{
		RawDiff: diff,
		Title:   pr.GetTitle(),
		Content: pr.GetBody(),
	}, nil
}
