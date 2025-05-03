package functions

const FuncSubmitFiles = "submit_files"

func InitSubmitFilesGitHubFunction(service SubmitFilesService) Function {
	f := Function{
		Name:        FuncSubmitFiles,
		Description: "Submit the modified files by Creation GitHub Pull Request",
		Func:        SubmitFileCaller(service),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"commit_message_short": map[string]any{
					"type":        "string",
					"description": "Short Commit message indicating purpose to change the file",
				},
				"commit_message_detail": map[string]any{
					"type":        "string",
					"description": "Detail commit message indicating changes to the file",
				},
				"pull_request_content": map[string]any{
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

func SubmitFileCaller(service SubmitFilesService) SubmitFilesType {
	return func(input SubmitFilesInput) (SubmitFilesOutput, error) {
		return service.SubmitFiles(input)
	}
}
