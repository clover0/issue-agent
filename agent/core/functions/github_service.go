package functions

type GitHubService interface {
	GetIssue(repository string, prNumber string) (GetIssueOutput, error)
	GetPullRequest(prNumber string) (GetPullRequestOutput, error)
	GetRepositoryContent(input GetRepositoryContentInput) (GetRepositoryContentOutput, error)

	CreateIssueComment(issueNumber string, comment string) (CreateIssueCommentOutput, error)
	CreateReviewCommentOne(input CreatePullRequestReviewCommentInput) (CreatePullRequestReviewCommentOutput, error)
}
