package functions

const FuncCreatePullRequestReviewComment = "create_pull_request_review_comment"

type CreatePullRequestReviewCommentType func(input CreatePullRequestReviewCommentInput) (CreatePullRequestReviewCommentOutput, error)

func InitCreatePullRequestReviewCommentFunction(service GitHubService) Function {
	f := Function{
		Name:        FuncCreatePullRequestReviewComment,
		Description: "Create a review comment on a GitHub pull request for a specific file and line range.",
		Func:        CreatePullRequestReviewCommentCaller(service),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"pr_number": map[string]any{
					"type":        "string",
					"description": "GitHub Pull Request Number to create comment to",
				},
				"review_file_path": map[string]any{
					"type":        "string",
					"description": "File path from repository root for review",
				},
				"review_start_line": map[string]any{
					"type":        "number",
					"description": "Review start line number on file",
					"minimum":     1,
				},
				"review_end_line": map[string]any{
					"type":        "number",
					"description": "Review end line number on file",
					"minimum":     1,
				},
				"review_comment": map[string]any{
					"type":        "string",
					"description": "Comment to be added to the pull request review",
				},
			},
			"required":             []string{"pr_number", "review_file_path", "review_start_line", "review_end_line", "review_comment"},
			"additionalProperties": false,
		},
	}

	register(f)

	return f
}

type CreatePullRequestReviewCommentInput struct {
	PRNumber        string `json:"pr_number"`
	ReviewFilePath  string `json:"review_file_path"`
	ReviewStartLine int    `json:"review_start_line"`
	ReviewEndLine   int    `json:"review_end_line"`
	ReviewComment   string `json:"review_comment"`
}

type CreatePullRequestReviewCommentOutput struct{}

func (g CreatePullRequestReviewCommentOutput) ToLLMString() string {
	return "success creating pull request review comment."
}

func CreatePullRequestReviewCommentCaller(service GitHubService) CreatePullRequestReviewCommentType {
	return func(input CreatePullRequestReviewCommentInput) (CreatePullRequestReviewCommentOutput, error) {
		return service.CreateReviewCommentOne(input)
	}
}
