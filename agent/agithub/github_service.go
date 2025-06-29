package agithub

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v71/github"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/util/pointer"
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
		PRNumber: prNumber,
		Head:     pr.GetHead().GetRef(),
		Base:     pr.GetBase().GetRef(),
		RawDiff:  diff,
		Title:    pr.GetTitle(),
		Content:  pr.GetBody(),
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

const contentSeparator = "---"

func (s GitHubService) GetRepositoryContent(input functions.GetRepositoryContentInput) (functions.GetRepositoryContentOutput, error) {
	c := context.Background()

	content, dirContent, _, err := s.client.Repositories.GetContents(c, s.owner, input.RepositoryName, input.Path, nil)
	if err != nil {
		return functions.GetRepositoryContentOutput{}, fmt.Errorf("failed to get repository content: %w", err)
	}

	var contents []*github.RepositoryContent

	if content != nil {
		contents = append(contents, content)
	}

	if len(dirContent) > 0 {
		contents = append(contents, dirContent...)
	}

	var contentStr string
	for _, cont := range contents {
		decoded, err := cont.GetContent()
		if err != nil {
			return functions.GetRepositoryContentOutput{}, fmt.Errorf("failed to decode content: %w", err)
		}
		contentStr += fmt.Sprintf("%s\nfile name: %s\nfile path: %s\n\n%s",
			contentSeparator, cont.GetName(), cont.GetPath(), decoded)
	}

	return functions.GetRepositoryContentOutput{
		Content: contentStr,
	}, nil
}

func (s GitHubService) CreateIssueComment(issueNumber string, comment string) (functions.CreateIssueCommentOutput, error) {
	c := context.Background()
	number, err := strconv.Atoi(issueNumber)
	if err != nil {
		return functions.CreateIssueCommentOutput{}, fmt.Errorf("failed to convert issue number to int: %w", err)
	}

	issueComment := &github.IssueComment{Body: &comment}
	_, _, err = s.client.Issues.CreateComment(c, s.owner, s.repository, number, issueComment)
	if err != nil {
		return functions.CreateIssueCommentOutput{}, fmt.Errorf("failed to create issue comment: %w", err)
	}

	return functions.CreateIssueCommentOutput{}, nil
}

func (s GitHubService) CreateReviewCommentOne(review functions.CreatePullRequestReviewCommentInput) (functions.CreatePullRequestReviewCommentOutput, error) {
	c := context.Background()

	prNumber, err := strconv.Atoi(review.PRNumber)
	if err != nil {
		return functions.CreatePullRequestReviewCommentOutput{}, fmt.Errorf("failed to convert prNumber to int: %w", err)
	}

	startLine := pointer.Ptr(review.ReviewStartLine)
	if review.ReviewStartLine == review.ReviewEndLine {
		startLine = nil
	}

	reviewComment := []*github.DraftReviewComment{{
		Path:      pointer.Ptr(review.ReviewFilePath),
		Body:      pointer.Ptr(review.ReviewComment),
		Side:      pointer.Ptr("RIGHT"),
		StartLine: startLine,
		Line:      pointer.Ptr(review.ReviewEndLine),
	}}

	_, _, err = s.client.PullRequests.CreateReview(c, s.owner, s.repository, prNumber, &github.PullRequestReviewRequest{
		Event:    pointer.Ptr("COMMENT"),
		Comments: reviewComment,
	})
	if err != nil {
		return functions.CreatePullRequestReviewCommentOutput{}, fmt.Errorf("failed to create review comment: %w", err)
	}

	return functions.CreatePullRequestReviewCommentOutput{}, nil
}

func (s GitHubService) RequestReviewers(prNumber int, reviewers []string, teamReviewers []string) (functions.RequestReviewersOutput, error) {
	c := context.Background()

	_, resp, err := s.client.PullRequests.RequestReviewers(
		c,
		s.owner,
		s.repository,
		prNumber,
		github.ReviewersRequest{
			Reviewers:     reviewers,
			TeamReviewers: teamReviewers,
		})
	if err == nil {
		return functions.RequestReviewersOutput{}, nil
	}

	// client error handling
	if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return functions.RequestReviewersOutput{}, fmt.Errorf("failed to read response body: %w", err)
		}
		defer func(b io.ReadCloser) {
			err := b.Close()
			if err != nil {
				s.logger.Error("failed to close response body: %s", err)
			}
		}(resp.Body)

		return functions.RequestReviewersOutput{},
			fmt.Errorf("client error requesting reviewers=%s, team_reviewers=%s to PR: %s", reviewers, teamReviewers, body)
	}

	return functions.RequestReviewersOutput{},
		fmt.Errorf("failed to request reviewers=%s, team_reviewers=%s to PR: %w", reviewers, teamReviewers, err)
}
