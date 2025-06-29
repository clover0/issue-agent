package core

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v71/github"

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
	ctx context.Context,
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
		lo.Error("failed to select forwarder: %s\n", err)
		return err
	}

	promptTemplate, err := coreprompt.LoadPrompt()
	if err != nil {
		return fmt.Errorf("failed to load prompt template: %w", err)
	}

	// check if the base branch exists
	ghService := agithub.NewGitHubService(conf.Agent.GitHub.Owner, workRepository, gh, lo)
	if _, err = ghService.GetBranch(baseBranch); err != nil {
		return err
	}

	submitService := agithub.NewSubmitFileGitHubService(
		lo, gh,
		functions.SubmitFilesServiceInput{
			GitHubOwner: conf.Agent.GitHub.Owner,
			Repository:  workRepository,
			BaseBranch:  baseBranch,
			GitEmail:    conf.Agent.Git.UserEmail,
			GitName:     conf.Agent.Git.UserName,
			PRLabels:    conf.Agent.GitHub.PRLabels,
		})
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
		return fmt.Errorf("failed to get issue: %w", err)
	}

	prompt, err := coreprompt.BuildRequirementPrompt(promptTemplate, conf.Language, baseBranch, issue)
	if err != nil {
		lo.Error("failed build requirement prompt: %s\n", err)
		return err
	}
	requirementAgent, err := RunAgent("requirementAgent",
		prompt, parameter, lo, &dataStore, llmForwarder, PlanTools())
	if err != nil {
		lo.Error("requirement agent failed: %s\n", err)
		return err
	}

	instruction := requirementAgent.LastHistory().RawContent
	prompt, err = coreprompt.BuildDeveloperPrompt(promptTemplate, conf.Language, baseBranch, issue, instruction)
	if err != nil {
		lo.Error("failed build developer prompt: %s\n", err)
		return err
	}

	if _, err := RunAgent("developerAgent", prompt, parameter, lo, &dataStore, llmForwarder, tools); err != nil {
		lo.Error("developer agent failed: %s\n", err)
		return err
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
	// TODO: move selection llm forwarder
	llmForwarder, err := selectForward(lo, conf.Agent.Model)
	if err != nil {
		lo.Error("failed to select forwarder: %s\n", err)
		return err
	}

	promptTemplate, err := coreprompt.LoadPrompt()
	if err != nil {
		return fmt.Errorf("failed to load prompt template: %w", err)
	}

	ghService := agithub.NewGitHubService(conf.Agent.GitHub.Owner, workRepository, gh, lo)

	submitFilesService := agithub.NopSubmitFileService{}
	submitRevisionService := agithub.NewSubmitRevisionGitHubService(lo, gh,
		functions.SubmitRevisionServiceInput{
			GitHubOwner: conf.Agent.GitHub.Owner,
			Repository:  workRepository,
			BaseBranch:  pr.Base,
			WorkBranch:  pr.Head,
			GitEmail:    conf.Agent.Git.UserEmail,
			GitName:     conf.Agent.Git.UserName,
		})

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
			CommentingTools(),
		))

	tools := slices.Concat(CommentingTools(), InvokeAgentTools())

	lo.Info("allowed functions: %s\n", strings.Join(util.Map(
		tools,
		func(e functions.Function) string { return e.Name.String() },
	), ","))
	lo.Info("agents will push to %s/%s branch %s\n", conf.Agent.GitHub.Owner, workRepository, pr.Head)

	dataStore := corestore.NewStore()

	prompt, err := coreprompt.BuildCommentReactorPrompt(promptTemplate, conf.Language, comment, pr)
	if err != nil {
		lo.Error("failed build commentReactor prompt: %s\n", err)
		return err
	}

	_, err = RunAgent("commentReactorAgent",
		prompt, parameter, lo, &dataStore, llmForwarder,
		tools,
	)
	if err != nil {
		lo.Error("developer agent failed: %s\n", err)
		return err
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
		lo.Error("requirement agent failed: %s\n", err)
		return &Agent{}, err
	}

	return ag, nil
}
