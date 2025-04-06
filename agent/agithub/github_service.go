package agithub

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v69/github"

	"github.com/clover0/issue-agent/core/functions"
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

func (s GitHubService) GetIssue(repository string, issueNumber string) (functions.GetIssueOutput, error) {
	number, err := strconv.Atoi(issueNumber)
	if err != nil {
		return functions.GetIssueOutput{}, fmt.Errorf("failed to convert issue number to int: %w", err)
	}

	c := context.Background()
	issue, _, err := s.client.Issues.Get(c, s.owner, repository, number)
	if err != nil {
		return functions.GetIssueOutput{}, fmt.Errorf("failed to get issue: %w", err)
	}

	return functions.GetIssueOutput{
		Path:    strconv.Itoa(issue.GetNumber()),
		Title:   issue.GetTitle(),
		Content: issue.GetBody(),
	}, nil
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
		Head:    pr.GetHead().GetRef(),
		Base:    pr.GetBase().GetRef(),
		RawDiff: diff,
		Title:   pr.GetTitle(),
		Content: pr.GetBody(),
	}, nil
}

func (s GitHubService) GetBranch(branchName string) (string, error) {
	c := context.Background()
	branch, resp, err := s.client.Repositories.GetBranch(c, s.owner, s.repository, branchName, 0)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return "", fmt.Errorf("branch %s not found : %w", branchName, err)
		}
		return "", fmt.Errorf("failed to get branch: %w", err)
	}

	return branch.GetName(), nil
}

func (s GitHubService) GetComment(commentNumber string) (functions.GetCommentOutput, error) {
	c := context.Background()
	number, err := strconv.ParseInt(commentNumber, 10, 64)
	if err != nil {
		return functions.GetCommentOutput{}, fmt.Errorf("failed to convert comment number to int %s", commentNumber)
	}

	comment, _, err := s.client.Issues.GetComment(c, s.owner, s.repository, number)
	if err != nil {
		return functions.GetCommentOutput{}, fmt.Errorf("failed to get comment: %w", err)
	}

	u, err := url.Parse(*comment.IssueURL)
	if err != nil {
		return functions.GetCommentOutput{}, fmt.Errorf("failed to parse issue url: %w", err)
	}

	parts := strings.Split(u.Path, "/")
	issueNumber := parts[len(parts)-1]

	return functions.GetCommentOutput{
		IssueNumber: issueNumber,
		Content:     comment.GetBody(),
	}, nil
}

func (s GitHubService) GetReviewComment(reviewID string) (functions.GetReviewOutput, error) {
	c := context.Background()
	id, err := strconv.ParseInt(reviewID, 10, 64)
	if err != nil {
		return functions.GetReviewOutput{}, fmt.Errorf("failed to convert review id to int %s", reviewID)
	}

	review, _, err := s.client.PullRequests.GetComment(c, s.owner, s.repository, id)
	if err != nil {
		return functions.GetReviewOutput{}, fmt.Errorf("failed to get review: %w", err)
	}

	u, err := url.Parse(review.GetPullRequestURL())
	if err != nil {
		return functions.GetReviewOutput{}, fmt.Errorf("failed to parse pull request url: %w", err)
	}

	parts := strings.Split(u.Path, "/")
	issueNumber := parts[len(parts)-1]

	startLine := review.GetOriginalStartLine()
	if startLine == 0 {
		startLine = review.GetOriginalLine()
	}
	return functions.GetReviewOutput{
		IssuesNumber: issueNumber,
		Path:         review.GetPath(),
		StartLine:    startLine,
		EndLine:      review.GetOriginalLine(),
		Content:      review.GetBody(),
	}, nil

}
