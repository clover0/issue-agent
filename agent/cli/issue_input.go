package cli

import (
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/clover0/issue-agent/config"
)

type IssueInputs struct {
	Common            *CommonInput
	GithubIssueNumber string
	WorkRepository    string `validate:"required"`
	BaseBranch        string `validate:"required"`
	FromFile          string
}

// TODO: refactor to be common flags
func IssueFlags() (*flag.FlagSet, *IssueInputs) {
	flagMapper := &IssueInputs{
		Common: &CommonInput{},
	}

	cmd := flag.NewFlagSet("issue", flag.ExitOnError)

	addCommonFlags(cmd, flagMapper.Common)

	cmd.StringVar(&flagMapper.WorkRepository, "work_repository", "", "Working repository to develop and create pull request")
	cmd.StringVar(&flagMapper.GithubIssueNumber, "github_issue_number", "", "GitHub issue number to solve")
	cmd.StringVar(&flagMapper.BaseBranch, "base_branch", "", "Base Branch for pull request")
	cmd.StringVar(&flagMapper.FromFile, "from_file", "", "Issue content from file path")

	return cmd, flagMapper
}

func ParseIssueInput(flags []string) (IssueInputs, error) {
	cmd, cliIn := IssueFlags()
	if err := cmd.Parse(flags); err != nil {
		return IssueInputs{}, fmt.Errorf("failed to parse input: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(cliIn); err != nil {
		errs := err.(validator.ValidationErrors)
		return IssueInputs{}, fmt.Errorf("validation failed: %w", errs)
	}

	if cliIn.GithubIssueNumber == "" && cliIn.FromFile == "" {
		return IssueInputs{}, fmt.Errorf("github_issue_number or from_file is required")
	}

	return *cliIn, nil
}

func MergeIssueInputConfig(conf config.Config, inputs IssueInputs) config.Config {
	if inputs.Common.LogLevel != "" {
		conf.LogLevel = inputs.Common.LogLevel
	}

	if inputs.Common.Language != "" {
		conf.Language = inputs.Common.Language
	}

	if inputs.Common.Model != "" {
		conf.Agent.Model = inputs.Common.Model
	}

	if inputs.Common.GitHubOwner != "" {
		conf.Agent.GitHub.Owner = inputs.Common.GitHubOwner
	}

	return conf
}
