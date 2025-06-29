package functions

const FuncRequestReviewers = "request_reviewers"

type RequestReviewersType func(input RequestReviewersInput) (RequestReviewersOutput, error)

func InitRequestReviewersFunction(service GitHubService) Function {
	f := Function{
		Name:        FuncRequestReviewers,
		Description: "Request reviewers for a GitHub pull request.",
		Func:        RequestReviewersCaller(service),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"pr_number": map[string]any{
					"type":        "number",
					"description": "GitHub Pull Request Number to request reviewers for.",
				},
				"member_reviewers": map[string]any{
					"type": "array",
					"items": map[string]string{
						"type": "string",
					},
					"description": "List of member `login`s to request on the pull request.",
				},
				"team_reviewers": map[string]any{
					"type": "array",
					"items": map[string]string{
						"type": "string",
					},
					"description": "List of team `slug`s to request on the pull request.",
				},
			},
			"required":             []string{"pr_number"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncRequestReviewers] = f

	return f
}

type RequestReviewersInput struct {
	PRNumber        int      `json:"pr_number"`
	MemberReviewers []string `json:"member_reviewers"`
	TeamReviewers   []string `json:"team_reviewers"`
}

type RequestReviewersOutput struct{}

func (g RequestReviewersOutput) ToLLMString() string {
	return "success requesting reviewers for pull request."
}

func RequestReviewersCaller(service GitHubService) RequestReviewersType {
	return func(input RequestReviewersInput) (RequestReviewersOutput, error) {
		return service.RequestReviewers(input.PRNumber, input.MemberReviewers, input.TeamReviewers)
	}
}
