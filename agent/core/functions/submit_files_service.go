package functions

type SubmitFilesServiceInput struct {
	GitHubOwner string
	Repository  string
	BaseBranch  string
	GitEmail    string
	GitName     string
	PRLabels    []string
}

type SubmitFilesType func(input SubmitFilesInput) (SubmitFilesOutput, error)

type SubmitFilesOutput struct {
	Message           string
	PushedBranch      string
	PullRequestNumber int
}

type SubmitFilesService interface {
	SubmitFiles(callerInput SubmitFilesInput) (SubmitFilesOutput, error)
}
