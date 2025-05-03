package functions

const FuncSubmitRevision = "submit_revision"

func InitSubmitRevisionFunction(service SubmitRevisionService) Function {
	f := Function{
		Name:        FuncSubmitRevision,
		Description: "Submit revision commits changed files using git add and git commit, finally git push on working branch.",
		Func:        SubmitRevisionCaller(service),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"commit_message_short": map[string]any{
					"type":        "string",
					"description": "Short commit message indicating purpose to resubmit",
				},
				"commit_message_detail": map[string]any{
					"type":        "string",
					"description": "Detail commit message indicating resubmitting content",
				},
			},
			"required":             []string{"commit_message_short"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncSubmitRevision] = f

	return f
}

type SubmitRevisionInput struct {
	CommitMessageShort  string `json:"commit_message_short"`
	CommitMessageDetail string `json:"commit_message_detail"`
}

func SubmitRevisionCaller(service SubmitRevisionService) SubmitRevisionType {
	return func(input SubmitRevisionInput) (SubmitRevisionOutput, error) {
		return service.SubmitRevision(input)
	}
}
