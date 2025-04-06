package functions

import "fmt"

const FuncGetIssue = "get_issue"

type GetIssueType func(input GetIssueInput) (GetIssueOutput, error)

func InitGetIssueFunction(service GitHubService) Function {
	f := Function{
		Name:        FuncGetIssue,
		Description: "Get a GitHub issue from organization(owner) passed as CLI input.",
		Func:        GetIssueCaller(service),
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repository_name": map[string]interface{}{
					"type":        "string",
					"description": "GitHub repository name to get the issue from. The `repo` part of the `owner/repo` format.",
				},
				"issue_number": map[string]interface{}{
					"type":        "string",
					"description": "GitHub Issue Number to get",
				},
			},
			"required":             []string{"repository_name", "issue_number"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncGetIssue] = f

	return f
}

type GetIssueInput struct {
	RepositoryName string `json:"repository_name"`
	IssueNumber    string `json:"issue_number"`
}

type GetIssueOutput struct {
	Path    string
	Title   string
	Content string
}

func (g GetIssueOutput) ToLLMString() string {
	s := fmt.Sprintf("# Issue Number\n%s\n\n", g.Path)
	s += fmt.Sprintf("# Title:\n%s\n\n", g.Title)
	s += fmt.Sprintf("# Content:\n%s\n", g.Content)
	return s
}

func GetIssueCaller(service GitHubService) GetIssueType {
	return func(input GetIssueInput) (GetIssueOutput, error) {
		return service.GetIssue(input.RepositoryName, input.IssueNumber)
	}
}
