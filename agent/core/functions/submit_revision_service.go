package functions

type SubmitRevisionServiceInput struct {
	GitHubOwner string
	Repository  string
	BaseBranch  string
	WorkBranch  string
	GitEmail    string
	GitName     string
}

type SubmitRevisionType func(input SubmitRevisionInput) (SubmitRevisionOutput, error)

type SubmitRevisionOutput struct {
	Message string
}

type SubmitRevisionService interface {
	SubmitRevision(callerInput SubmitRevisionInput) (SubmitRevisionOutput, error)
}
