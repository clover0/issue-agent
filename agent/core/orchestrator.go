package core

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v73/github"

	"github.com/clover0/issue-agent/agithub"
	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/core/functions"
	coreprompt "github.com/clover0/issue-agent/core/prompt"
	corestore "github.com/clover0/issue-agent/core/store"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/util"
)

// OrchestrateAgentsByIssue orchestrates agents
// Currently, the processing is based on the issue command
// TODO: refactor many arguments
// TODO: no dependent on issue command
func OrchestrateAgentsByIssue(
	_ context.Context,
	lo logger.Logger,
	conf config.Config,
	baseBranch string,
	workRepository string,
	gh *github.Client,
	issueNumber string,
	selectForward SelectForwarder,
) error {
	llmForwarder, err := selectForward(lo, conf.Agent.Model)
	if err != nil {
		return fmt.Errorf("select forwarder: %w", err)
	}

	// check if the base branch exists
	ghService := agithub.NewGitHubService(conf.Agent.GitHub.Owner, workRepository, gh, lo)
	if _, err = ghService.GetBranch(baseBranch); err != nil {
		return err
	}

	submitService, err := agithub.NewSubmitFileGitHubService(
		lo, gh,
		functions.SubmitFilesServiceInput{
			GitHubOwner: conf.Agent.GitHub.Owner,
			Repository:  workRepository,
			BaseBranch:  baseBranch,
			GitEmail:    conf.Agent.Git.UserEmail,
			GitName:     conf.Agent.Git.UserName,
			PRLabels:    conf.Agent.GitHub.PRLabels,
		})
	if err != nil {
		return fmt.Errorf("create submit file service: %w", err)
	}

	submitRevisionService := agithub.NopSubmitRevisionService{}

	parameter := Parameter{
		MaxSteps: conf.Agent.MaxSteps,
		Model:    conf.Agent.Model,
	}

	functions.InitializeFunctions(
		*conf.Agent.GitHub.NoSubmit,
		ghService,
		submitService,
		submitRevisionService,
		conf.Agent.AllowFunctions,
	)

	tools := functions.AllFunctions()
	functions.InitializeInvokeAgentFunction(
		conf.Agent.AllowFunctions,
		NewAgentInvoker(
			parameter,
			lo,
			llmForwarder,
			tools,
		))

	tools = functions.AllFunctions()
	lo.Info("allowed functions: %s\n", strings.Join(util.Map(
		tools,
		func(e functions.Function) string { return e.Name.String() },
	), ","))
	lo.Info("agents make a pull request to %s/%s\n", conf.Agent.GitHub.Owner, workRepository)

	dataStore := corestore.NewStore()

	issue, err := ghService.GetIssue(workRepository, issueNumber)
	if err != nil {
		return fmt.Errorf("get issue: %w", err)
	}

	prompt, err := coreprompt.PlanningPrompt{
		Language:     conf.Language,
		BaseBranch:   baseBranch,
		IssueTitle:   issue.Title,
		IssueContent: issue.Content,
		IssueNumber:  issue.Path,
	}.Build()
	if err != nil {
		return fmt.Errorf("orchestrator builds planning prompt: %w", err)
	}
	planningAgent, err := RunAgent("planningAgent",
		prompt, parameter, lo, &dataStore, llmForwarder, PlanTools())
	if err != nil {
		return err
	}

	instruction := planningAgent.LastHistory().RawContent
	prompt, err = coreprompt.DeveloperPrompt{
		Language:     conf.Language,
		BaseBranch:   baseBranch,
		IssueTitle:   issue.Title,
		IssueContent: issue.Content,
		IssueNumber:  issue.Path,
		Instruction:  instruction,
	}.Build()
	if err != nil {
		return fmt.Errorf("orchestrator builds developer prompt: %w", err)
	}

	if _, err := RunAgent("developerAgent", prompt, parameter, lo, &dataStore, llmForwarder, tools); err != nil {
		return fmt.Errorf("orchestrator developer agent: %w", err)
	}

	lo.Info("agents finished work\n")

	return nil
}

func OrchestrateAgentsByComment(
	lo logger.Logger,
	conf config.Config,
	workRepository string,
	gh *github.Client,
	selectForward SelectForwarder,
	comment functions.GetCommentOutput,
	pr functions.GetPullRequestOutput,
) error {
	llmForwarder, err := selectForward(lo, conf.Agent.Model)
	if err != nil {
		return fmt.Errorf("select forwarder: %w", err)
	}

	ghService := agithub.NewGitHubService(conf.Agent.GitHub.Owner, workRepository, gh, lo)

	submitFilesService := agithub.NopSubmitFileService{}
	submitRevisionService, err := agithub.NewSubmitRevisionGitHubService(lo, gh,
		functions.SubmitRevisionServiceInput{
			GitHubOwner: conf.Agent.GitHub.Owner,
			Repository:  workRepository,
			BaseBranch:  pr.Base,
			WorkBranch:  pr.Head,
			GitEmail:    conf.Agent.Git.UserEmail,
			GitName:     conf.Agent.Git.UserName,
		})
	if err != nil {
		return fmt.Errorf("create submit revision service: %w", err)
	}

	parameter := Parameter{
		MaxSteps: conf.Agent.MaxSteps,
		Model:    conf.Agent.Model,
	}

	functions.InitializeFunctions(
		*conf.Agent.GitHub.NoSubmit,
		ghService,
		submitFilesService,
		submitRevisionService,
		conf.Agent.AllowFunctions,
	)

	functions.InitializeInvokeAgentFunction(
		conf.Agent.AllowFunctions,
		NewAgentInvoker(
			parameter,
			lo,
			llmForwarder,
			ReactTools(),
		))

	tools := slices.Concat(ReactTools(), InvokeAgentTools())

	lo.Info("allowed functions: %s\n", strings.Join(util.Map(
		tools,
		func(e functions.Function) string { return e.Name.String() },
	), ","))
	lo.Info("agents will push to %s/%s branch %s\n", conf.Agent.GitHub.Owner, workRepository, pr.Head)

	dataStore := corestore.NewStore()

	prompt, err := coreprompt.CommentReactorPrompt{
		Language:      conf.Language,
		WorkingBranch: pr.Head,
		PRNumber:      pr.PRNumber,
		Comment:       comment.Content,
		PRLLMString:   pr.ToLLMString(),
	}.Build()
	if err != nil {
		lo.Error("orchestrator builds comment reactor prompt: %s\n", err)
		return err
	}

	_, err = RunAgent("commentReactorAgent",
		prompt, parameter, lo, &dataStore, llmForwarder,
		tools,
	)
	if err != nil {
		return fmt.Errorf("orchestrator comment reactor agent: %w", err)
	}
	lo.Info("agents finished work\n")

	return nil
}

func RunAgent(
	name string,
	prompt coreprompt.Prompt,
	parameter Parameter,
	lo logger.Logger,
	dataStore *corestore.Store,
	llmForwarder LLMForwarder,
	tools []functions.Function,
) (AgentLike, error) {
	ag := NewAgent(
		parameter,
		name,
		lo,
		prompt,
		llmForwarder,
		dataStore,
		tools,
	)

	if _, err := ag.Work(); err != nil {
		return &Agent{}, fmt.Errorf("agent %s failed: %w", name, err)
	}

	return ag, nil
}
