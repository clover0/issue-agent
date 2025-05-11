package functions

import (
	"bytes"
	"fmt"
	"text/template"
)

const FuncGetPullRequest = "get_pull_request"

// TODO: move to a separate file
type GitHubService interface {
	GetIssue(repository string, prNumber string) (GetIssueOutput, error)
	GetPullRequest(prNumber string) (GetPullRequestOutput, error)
	GetRepositoryContent(input GetRepositoryContentInput) (GetRepositoryContentOutput, error)

	CreateIssueComment(issueNumber string, comment string) (CreateIssueCommentOutput, error)
	CreateReviewCommentOne(input CreatePullRequestReviewCommentInput) (CreatePullRequestReviewCommentOutput, error)
}

type GetPullRequestType func(input GetPullRequestInput) (GetPullRequestOutput, error)

func InitGetPullRequestFunction(service GitHubService) Function {
	f := Function{
		Name:        FuncGetPullRequest,
		Description: "Get a GitHub Pull Request",
		Func:        GetPullRequestCaller(service),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"pr_number": map[string]any{
					"type":        "string",
					"description": "Pull Request Number to get",
				},
			},
			"required":             []string{"pr_number"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncGetPullRequest] = f

	return f
}

type GetPullRequestInput struct {
	PRNumber string `json:"pr_number"`
}

type GetPullRequestOutput struct {
	PRNumber string
	Head     string
	Base     string
	RawDiff  string
	Title    string
	Content  string
}

type GetCommentOutput struct {
	// IssueNumber is the issue number.
	// When comment is on the pull request, IssueNumber is the pull request number.
	IssueNumber string

	Content string
}

type GetReviewOutput struct {
	IssuesNumber string
	Path         string
	StartLine    int
	EndLine      int
	Content      string
}

func (g GetReviewOutput) ToLLMString() string {
	tmpl := `
The following file information received a code review.

# Review information
* Review file path: {{ .Path }}
* Review start line number: {{ .StartLine }}
* Review end line number: {{ .EndLine }}

# Review content
{{ .Content }}
`
	t, err := template.New("reviewComments").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse review template: %w", err).Error()
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, g)
	if err != nil {
		return fmt.Errorf("failed to execute review template: %w", err).Error()
	}

	return buf.String()
}

func (g GetPullRequestOutput) ToLLMString() string {
	errMsg := "failed to convert pull-request to string for LLM"

	tmpl := `
<pr-number>
{{ .PRNumber }}
</pr-number>

<pull-request-title>
{{ .Title }}
</pull-request-title>

<pull-request-description>
{{ .Content }}
</pull-request-description>

<pull-request-diff>
{{ .RawDiff }}
</pull-request-diff>
`

	t, err := template.New("pullRequest").Parse(tmpl)
	if err != nil {
		return errMsg
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, g)
	if err != nil {
		return errMsg
	}

	return buf.String()
}

func GetPullRequestCaller(service GitHubService) GetPullRequestType {
	return func(input GetPullRequestInput) (GetPullRequestOutput, error) {
		return service.GetPullRequest(input.PRNumber)
	}
}
