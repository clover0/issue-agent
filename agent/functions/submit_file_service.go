package functions

import "context"

type SubmitFilesServiceInput struct {
	BaseBranch string
	GitEmail   string
	GitName    string
	PRLabels   []string
}

type SubmitFilesCallerType func(input SubmitFilesInput) (SubmitFilesOutput, error)

type SubmitFilesOutput struct {
	Message           string
	PushedBranch      string
	PullRequestNumber int
}

type SubmitFilesService interface {
	Caller(ctx context.Context, input SubmitFilesServiceInput) SubmitFilesCallerType
}
