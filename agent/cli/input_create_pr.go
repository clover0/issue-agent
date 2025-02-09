package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/clover0/issue-agent/config"
)

type CreatePRInput struct {
	Common            *CommonInput
	GitHubOwner       string `validate:"required"`
	GithubIssueNumber string
	WorkRepository    string `validate:"required"`
	BaseBranch        string `validate:"required"`
	FromFile          string
}

func (c *CreatePRInput) MergeGitHubArg(pr ArgGitHubCreatePR) *CreatePRInput {
	c.GitHubOwner = pr.Owner
	c.WorkRepository = pr.Repository
	c.GithubIssueNumber = pr.IssueNumber
	return c
}

func (c *CreatePRInput) MergeConfig(conf config.Config) config.Config {
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

	if c.Common.ReviewAgents > 0 {
		conf.Agent.ReviewAgents = c.Common.ReviewAgents
	}

	return conf
}

// Validate
func (c *CreatePRInput) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		// TODO: error message
		errs := err.(validator.ValidationErrors)
		return fmt.Errorf("validation failed: %w", errs)
	}

	if c.GithubIssueNumber == "" && c.FromFile == "" {
		return fmt.Errorf("github_issue_number or from_file is required")
	}

	return nil
}

func CreatePRFlags() (*flag.FlagSet, *CreatePRInput) {
	flagMapper := &CreatePRInput{
		Common: &CommonInput{},
	}

	cmd := flag.NewFlagSet("issue", flag.ExitOnError)

	addCommonFlags(cmd, flagMapper.Common)

	cmd.StringVar(&flagMapper.BaseBranch, "base_branch", "", "Base Branch for pull request")

	// NOTE: from_file is not implemented yet.
	// cmd.StringVar(&flagMapper.FromFile, "from_file", "", "Issue content from file path.")

	return cmd, flagMapper
}

type ArgGitHubCreatePR struct {
	Owner       string
	Repository  string
	IssueNumber string
}

// ParseGitHubArg binds the input to the GitHub input
// expected format: OWNER/REPOSITORY/issues/NUMBER
func ParseGitHubArg(arg string) (ArgGitHubCreatePR, error) {
	splits := strings.Split(arg, "/")
	if len(splits) != 4 {
		return ArgGitHubCreatePR{}, fmt.Errorf("invalid input format: `%s`. valid format is `OWNER/REPOSITORY/issues/NUMBER`", arg)
	}

	return ArgGitHubCreatePR{
		Owner:       splits[0],
		Repository:  splits[1],
		IssueNumber: splits[3],
	}, nil

}

func ParseCreatePRInput(argAndFlags []string) (CreatePRInput, error) {
	arg, flags := ParseArgFlags(argAndFlags)
	ghIn, err := ParseGitHubArg(arg)
	if err != nil {
		return CreatePRInput{}, fmt.Errorf("failed to parse arg: %w", err)
	}

	cmd, cliIn := CreatePRFlags()
	if err := cmd.Parse(flags); err != nil {
		return CreatePRInput{}, fmt.Errorf("failed to parse input: %w", err)
	}

	cliIn.MergeGitHubArg(ghIn)

	if err := cliIn.Validate(); err != nil {
		return CreatePRInput{}, err
	}

	return *cliIn, nil
}
