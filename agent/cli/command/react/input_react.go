package react

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"

	"github.com/clover0/issue-agent/cli/command/common"
	"github.com/clover0/issue-agent/cli/util"
	"github.com/clover0/issue-agent/config"
)

type ReactType string

const (
	Comment       ReactType = ReactType("comment")
	ReviewComment ReactType = ReactType("review_comment")
)

type ArgGitHubReact struct {
	ReactType ReactType

	Owner      string
	Repository string
	PRNumber   string
	CommentID  string
	ReviewID   string
}

type ReactInput struct {
	ReactType ReactType

	Common         *common.CommonInput
	GitHubOwner    string `validate:"required"`
	GithubPRNumber string
	WorkRepository string `validate:"required"`
	CommentID      string
	ReviewID       string
}

func (c *ReactInput) MergeGitHubArg(react ArgGitHubReact) *ReactInput {
	c.ReactType = react.ReactType
	c.GitHubOwner = react.Owner
	c.GithubPRNumber = react.PRNumber
	c.WorkRepository = react.Repository
	c.CommentID = react.CommentID
	c.ReviewID = react.ReviewID

	return c
}

func (c *ReactInput) MergeConfig(conf config.Config) config.Config {
	if c.Common.LogLevel != "" {
		conf.LogLevel = c.Common.LogLevel
	}

	if c.Common.Language != "" {
		conf.Language = c.Common.Language
	}

	if c.Common.Model != "" {
		conf.Agent.Model = c.Common.Model
	}

	if c.GitHubOwner != "" {
		conf.Agent.GitHub.Owner = c.GitHubOwner
	}

	return conf
}

func (c *ReactInput) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		// TODO: error message
		errs := err.(validator.ValidationErrors)
		return fmt.Errorf("validation failed: %w", errs)
	}

	return nil
}

// BindReactGitHubArg binds the input to the GitHub input
// Expect the input to be in the bellow format.
//
// issue_comment(pull request): OWNER/REPO/issues/comments/COMMENT_ID
// pull_request_review_comment: OWNER/REPO/pulls/comments/COMMENT_ID
func BindReactGitHubArg(arg string) (ArgGitHubReact, error) {
	commonPattern := `^(?P<owner>[^/]+)/(?P<repo>[^/]+)/`
	issueCommentPattern := commonPattern + `issues/comments/(?P<commentID>[^/]+)$`
	pullRequestReviewPattern := commonPattern + `pulls/comments/(?P<commentID>[^/]+)$`

	{
		// handle pull request review comment
		re := regexp.MustCompile(pullRequestReviewPattern)
		matches := re.FindStringSubmatch(arg)
		if len(matches) == 1+3 {
			return ArgGitHubReact{
				ReactType:  ReviewComment,
				Owner:      matches[re.SubexpIndex("owner")],
				Repository: matches[re.SubexpIndex("repo")],
				PRNumber:   "",
				ReviewID:   matches[re.SubexpIndex("commentID")],
			}, nil
		}
	}

	{
		// handle issue comment
		re := regexp.MustCompile(issueCommentPattern)
		matches := re.FindStringSubmatch(arg)
		if len(matches) == 1+3 {
			return ArgGitHubReact{
				ReactType:  Comment,
				Owner:      matches[re.SubexpIndex("owner")],
				Repository: matches[re.SubexpIndex("repo")],
				PRNumber:   "",
				CommentID:  matches[re.SubexpIndex("commentID")],
			}, nil
		}
	}

	return ArgGitHubReact{}, fmt.Errorf("failed to parse github arg: %s", arg)

}

func ReactFlags() (*flag.FlagSet, *ReactInput) {
	flagMapper := &ReactInput{
		Common: &common.CommonInput{},
	}

	cmd := flag.NewFlagSet("react", flag.ExitOnError)

	common.AddCommonFlags(cmd, flagMapper.Common)

	return cmd, flagMapper
}

func ParseReactInput(argAndFlags []string) (ReactInput, error) {
	arg, flags := util.ParseArgFlags(argAndFlags)
	ghIn, err := BindReactGitHubArg(arg)
	if err != nil {
		return ReactInput{}, fmt.Errorf("failed to parse arg: %w", err)
	}

	cmd, cliIn := ReactFlags()
	if err := cmd.Parse(flags); err != nil {
		return ReactInput{}, fmt.Errorf("failed to parse input: %w", err)
	}

	cliIn.MergeGitHubArg(ghIn)

	if err := cliIn.Validate(); err != nil {
		return ReactInput{}, err
	}

	return *cliIn, nil
}
