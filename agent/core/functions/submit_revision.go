package functions

const FuncSubmitRevision = "submit_revision"

func InitSubmitRevisionFunction(service SubmitRevisionService) Function {
	f := Function{
		Name:        FuncSubmitRevision,
		Description: "Submit revision is a function to correct and resubmit after submission.",
		Func:        SubmitRevisionCaller(service),
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"commit_message_short": map[string]interface{}{
					"type":        "string",
					"description": "Short commit message indicating purpose to resubmit",
				},
				"commit_message_detail": map[string]interface{}{
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
