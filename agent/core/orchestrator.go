package core

import (
	"bytes"
	"context"
	"encoding/json"
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
	"github.com/clover0/issue-agent/util/pointer"
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
		return fmt.Errorf("failed to load prompt template: %w\n", err)
	}

	// check if the base branch exists
	ghService := agithub.NewGitHubService(conf.Agent.GitHub.Owner, workRepository, gh, lo)
	if _, err = ghService.GetBranch(baseBranch); err != nil {
		return err
	}

	submitService := agithub.NewSubmitFileGitHubService(
		lo, gh,
		functions.SubmitFilesServiceInput{
			GitHubOwner:   conf.Agent.GitHub.Owner,
			Repository:    workRepository,
			BaseBranch:    baseBranch,
			GitEmail:      conf.Agent.Git.UserEmail,
			GitName:       conf.Agent.Git.UserName,
			PRLabels:      conf.Agent.GitHub.PRLabels,
			Reviewers:     conf.Agent.GitHub.Reviewers,
			TeamReviewers: conf.Agent.GitHub.TeamReviewers,
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
	developerAgent, err := RunAgent("developerAgent",
		prompt, parameter, lo, &dataStore, llmForwarder, tools)
	if err != nil {
		lo.Error("developer agent failed: %s\n", err)
		return err
	}

	if conf.Agent.ReviewAgents == 0 {
		lo.Info("skip review agents\n")
		lo.Info("agents finished work\n")
		return nil
	}

	if s := dataStore.GetSubmittedWork(corestore.LastSubmissionKey); s == nil {
		lo.Error("submission is not found\n")
		return err
	}
	submittedPRNumber := dataStore.GetSubmittedWork(corestore.LastSubmissionKey).PullRequestNumber

	prompt, err = coreprompt.BuildReviewManagerPrompt(
		promptTemplate, conf, issue, util.Map(developerAgent.ChangedFiles(), func(f corestore.File) string { return f.Path }), baseBranch)
	if err != nil {
		lo.Error("failed to build review manager prompt: %s\n", err)
		return err
	}
	reviewManager, err := RunAgent(
		"reviewManagerAgent",
		prompt,
		parameter,
		lo, &dataStore, llmForwarder, functions.AllFunctions())
	if err != nil {
		lo.Error("reviewManagerAgent failed: %s\n", err)
		return err
	}
	output := reviewManager.LastHistory().RawContent
	lo.Info("ReviewManagerAgent output: %s\n", output)
	type agentPrompt struct {
		AgentName string `json:"agent_name"`
		Prompt    string `json:"prompt"`
	}

	// FIXME: these code is fragile
	// parse json output for reviewer agents
	// expected output:
	//   text text text...
	//   [{"agent_name": "agent1", "prompt": "prompt1"}, ...]
	//   test...
	var prompts []agentPrompt
	jsonStart := strings.Index(output, "[")   // find JSON start
	jsonEnd := strings.LastIndex(output, "]") // find JSON end
	jsonLike := output[jsonStart : jsonEnd+1]
	jsonLike = strings.ReplaceAll(jsonLike, "\n", `\\n`)
	outBuff := bytes.NewBufferString(jsonLike)
	if err := json.Unmarshal(outBuff.Bytes(), &prompts); err != nil {
		lo.Error("failed to unmarshal output: %s\n", err)
		return err
	}

	for _, p := range prompts {
		lo.Info("Run %s\n", p.AgentName)
		prpt, err := coreprompt.BuildReviewerPrompt(promptTemplate, conf.Language, submittedPRNumber, p.Prompt)
		if err != nil {
			lo.Error("failed to build reviewer prompt: %s\n", err)
			return err
		}

		reviewer, err := RunAgent(
			p.AgentName,
			prpt,
			parameter,
			lo, &dataStore, llmForwarder, functions.AllFunctions())
		if err != nil {
			lo.Error("%s failed: %s\n", p.AgentName, err)
			return err
		}
		output := reviewer.LastHistory().RawContent

		// parse JSON output
		// TODO: validate
		var reviews []struct {
			ReviewFilePath  string `json:"review_file_path"`
			ReviewStartLine int    `json:"review_start_line"`
			ReviewEndLine   int    `json:"review_end_line"`
			ReviewComment   string `json:"review_comment"`
			Suggestion      string `json:"suggestion"`
		}
		jsonStart := strings.Index(output, "[")   // find JSON start
		jsonEnd := strings.LastIndex(output, "]") // find JSON end
		jsonStr := strings.ReplaceAll(output[jsonStart:jsonEnd+1], "\n", `\\n`)
		outBuff := bytes.NewBufferString(jsonStr)
		if err := json.Unmarshal(outBuff.Bytes(), &reviews); err != nil {
			lo.Error("failed to unmarshal output: %s\n", err)
			return err
		}

		// TODO: move to agithub package
		var comments []*github.DraftReviewComment
		for _, r := range reviews {
			startLine := pointer.Ptr(r.ReviewStartLine)
			if *startLine == 0 {
				*startLine = 1
			}
			if r.ReviewStartLine == r.ReviewEndLine {
				startLine = nil
			}
			body := fmt.Sprintf("from %s\n", p.AgentName) +
				r.ReviewComment
			if r.Suggestion != "" {
				// TODO: escape JSON in Suggestion string
				body += "\n\n```suggestion\n" + r.Suggestion + "\n```\n"
			}
			comments = append(comments, &github.DraftReviewComment{
				Path:      pointer.Ptr(r.ReviewFilePath),
				Body:      pointer.Ptr(body),
				StartLine: startLine,
				Line:      pointer.Ptr(r.ReviewEndLine),
				Side:      pointer.Ptr("RIGHT"),
			})
		}

		if _, _, err := gh.PullRequests.CreateReview(context.Background(),
			conf.Agent.GitHub.Owner,
			workRepository,
			submittedPRNumber,
			&github.PullRequestReviewRequest{
				Event:    pointer.Ptr("COMMENT"),
				Comments: comments,
			},
		); err != nil {
			lo.Error("failed to create pull request review: %s. but agent continue to work\n", err)
		}
		lo.Info("finish %s\n", p.AgentName)
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
		return fmt.Errorf("failed to load prompt template: %w\n", err)
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
