package functions

const FuncCreatePullRequestComment = "create_pull_request_comment"

type CreatePullRequestCommentType func(input CreatePullRequestCommentInput) (CreateIssueCommentOutput, error)

func InitCreatePullRequestCommentFunction(service GitHubService) Function {
	f := Function{
		Name:        FuncCreatePullRequestComment,
		Description: "Create a comment on a GitHub pull request from `owner/repo` passed as CLI input.",
		Func:        CreatePullRequestCommentCaller(service),
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pr_number": map[string]interface{}{
					"type":        "string",
					"description": "GitHub Pull Request Number to create comment to",
				},
				"comment": map[string]interface{}{
					"type":        "string",
					"description": "Comment by markdown on the pull request",
				},
			},
			"required":             []string{"pr_number", "comment"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncCreatePullRequestComment] = f

	return f
}

type CreatePullRequestCommentInput struct {
	PRNumber string `json:"pr_number"`
	Comment  string `json:"comment"`
}

type CreateIssueCommentOutput struct{}

func (g CreateIssueCommentOutput) ToLLMString() string {
	return "success creating pull request comment."
}

func CreatePullRequestCommentCaller(service GitHubService) CreatePullRequestCommentType {
	return func(input CreatePullRequestCommentInput) (CreateIssueCommentOutput, error) {
		return service.CreateIssueComment(input.PRNumber, input.Comment)
	}
}
