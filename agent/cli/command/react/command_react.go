package react

import (
	"fmt"
	"os"

	"github.com/clover0/issue-agent/agithub"
	"github.com/clover0/issue-agent/cli/util"
	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/core"
	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/models"
)

const ReactCommand = "react"

func React(flags []string) error {
	cliIn, err := ParseReactInput(flags)
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	conf, err := config.LoadDefault(util.IsPassedConfig(cliIn.Common.Config))
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	conf = cliIn.MergeConfig(conf)

	if err := config.ValidateConfig(conf); err != nil {
		return err
	}

	lo := logger.NewPrinter(conf.LogLevel)

	gh, err := agithub.NewGitHub()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	ghService := agithub.NewGitHubService(conf.Agent.GitHub.Owner, cliIn.WorkRepository, gh, lo)

	comment, err := getComment(ghService, cliIn)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}
	pr, err := ghService.GetPullRequest(comment.IssueNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}
	if pr.Base == pr.Head {
		lo.Info(fmt.Sprintf("base and head are the same. base=%s, head=%s\n", pr.Base, pr.Head))
		return nil
	}

	if *conf.Agent.GitHub.CloneRepository {
		if err := agithub.CloneRepository(lo, conf, cliIn.WorkRepository, pr.Head); err != nil {
			lo.Error("failed to clone repository")
			return err
		}
	}

	// TODO: no dependency with changing directory
	if err := os.Chdir(conf.WorkDir); err != nil {
		lo.Error("failed to change directory: %s\n", err)
		return err
	}

	return core.OrchestrateAgentsByComment(
		lo, conf, cliIn.WorkRepository, gh, models.SelectForwarder, comment, pr)
}

func getComment(ghService agithub.GitHubService, in ReactInput) (functions.GetCommentOutput, error) {
	switch in.ReactType {
	case Comment:
		comment, err := ghService.GetComment(in.CommentID)
		if err != nil {
			return functions.GetCommentOutput{}, fmt.Errorf("failed to get issue: %w", err)
		}
		return comment, nil

	case ReviewComment:
		comment, err := ghService.GetReviewComment(in.ReviewID)
		if err != nil {
			return functions.GetCommentOutput{}, fmt.Errorf("failed to get review: %w", err)
		}
		return functions.GetCommentOutput{
			IssueNumber: comment.IssuesNumber,
			Content:     comment.ToLLMString(),
		}, nil
	}

	return functions.GetCommentOutput{}, fmt.Errorf("invalid react type: %v", in.ReactType)
}
