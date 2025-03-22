package functions

import (
	"bytes"
	"fmt"
	"text/template"
)

const FuncGetPullRequest = "get_pull_request"

type RepositoryService interface {
	GetPullRequest(prNumber string) (GetPullRequestOutput, error)
}

type GetPullRequestType func(input GetPullRequestInput) (GetPullRequestOutput, error)

func InitGetPullRequestFunction(service RepositoryService) Function {
	f := Function{
		Name:        FuncGetPullRequest,
		Description: "Get a GitHub Pull Request",
		Func:        GetPullRequestCaller(service),
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pr_number": map[string]interface{}{
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
	Head    string
	Base    string
	RawDiff string
	Title   string
	Content string
}

type GetIssueOutput struct {
	Path    string
	Content string
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
	DiffHunk     string
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

# Review diff hunk
{{ .DiffHunk }}
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

func GetPullRequestCaller(service RepositoryService) GetPullRequestType {
	return func(input GetPullRequestInput) (GetPullRequestOutput, error) {
		return service.GetPullRequest(input.PRNumber)
	}
}
