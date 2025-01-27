package loader

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v68/github"
)

type GitHubLoader struct {
	client     *github.Client
	owner      string
	repository string
}

func NewGitHubLoader(gh *github.Client, owner, repository string) Loader {
	return GitHubLoader{
		client:     gh,
		owner:      owner,
		repository: repository,
	}
}

func (g GitHubLoader) LoadIssue(ctx context.Context, number string) (Issue, error) {
	num, err := strconv.Atoi(number)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to convert issue number to int: %w", err)
	}
	issue, _, err := g.client.Issues.Get(ctx, g.owner, g.repository, num)
	if err != nil {
		return Issue{}, err
	}

	return Issue{
		Content: *issue.Body,
		Path:    number,
	}, nil
}
