package functions

const FuncSubmitFiles = "submit_files"

func InitSubmitFilesGitHubFunction() Function {
	f := Function{
		Name:        FuncSubmitFiles,
		Description: "Submit the modified files by Creation GitHub Pull Request",
		Func:        SubmitFiles,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"commit_message_short": map[string]interface{}{
					"type":        "string",
					"description": "Short Commit message indicating purpose to change the file",
				},
				"commit_message_detail": map[string]interface{}{
					"type":        "string",
					"description": "Detail commit message indicating changes to the file",
				},
				"pull_request_content": map[string]interface{}{
					"type":        "string",
					"description": "Pull Request Content",
				},
			},
			"required":             []string{"commit_message_short", "pull_request_content"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncSubmitFiles] = f

	return f
}

type SubmitFilesInput struct {
	CommitMessageShort  string `json:"commit_message_short"`
	CommitMessageDetail string `json:"commit_message_detail"`
	PullRequestContent  string `json:"pull_request_content"`
}

func SubmitFiles(submitting SubmitFilesCallerType, input SubmitFilesInput) (SubmitFilesOutput, error) {
	return submitting(input)
}
