package functions

const FuncGetRepositoryContent = "get_repository_content"

type GetRepositoryContentType func(input GetRepositoryContentInput) (GetRepositoryContentOutput, error)

func InitGetRepositoryContentFunction(service GitHubService) Function {
	f := Function{
		Name:        FuncGetRepositoryContent,
		Description: "Get contents of a file or directory in a GitHub repository.",
		Func:        GetRepositoryContentCaller(service),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"repository_name": map[string]any{
					"type":        "string",
					"description": "GitHub repository name to get the content. This is `repo` part of the `owner/repo` format.",
				},
				"path": map[string]any{
					"type":        "string",
					"description": "File path from repository root.",
				},
			},
			"required":             []string{"repository_name", "path"},
			"additionalProperties": false,
		},
	}

	register(f)

	return f
}

type GetRepositoryContentInput struct {
	RepositoryName string `json:"repository_name"`
	Path           string `json:"path"`
}

type GetRepositoryContentOutput struct {
	Content string
}

func (g GetRepositoryContentOutput) ToLLMString() string {
	return g.Content
}

func GetRepositoryContentCaller(service GitHubService) GetRepositoryContentType {
	return func(input GetRepositoryContentInput) (GetRepositoryContentOutput, error) {
		return service.GetRepositoryContent(input)
	}
}
