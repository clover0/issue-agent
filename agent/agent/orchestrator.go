package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v68/github"

	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/functions"
	"github.com/clover0/issue-agent/functions/agithub"
	"github.com/clover0/issue-agent/loader"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/models"
	libprompt "github.com/clover0/issue-agent/prompt"
	"github.com/clover0/issue-agent/store"
	"github.com/clover0/issue-agent/util"
	"github.com/clover0/issue-agent/util/pointer"
)

// OrchestrateAgents orchestrates agents
// Currently, the processing is based on the issue command
// TODO: refactor many arguments
// TODO: no dependent on issue command
func OrchestrateAgents(
	ctx context.Context,
	lo logger.Logger,
	conf config.Config,
	loaderr loader.Loader,
	baseBranch string,
	issue loader.Issue,
	workRepository string,
	gh *github.Client,
) error {
	llmForwarder, err := models.SelectForwarder(lo, conf.Agent.Model)
	if err != nil {
		lo.Error("failed to select forwarder: %s\n", err)
		return err
	}

	promptPath := conf.Agent.PromptPath
	if len(promptPath) > 0 {
		// In container, the prompt file is mounted to `config.PromptFilePath`
		promptPath = config.PromptFilePath
	}
	promptTemplate, err := libprompt.LoadPrompt(promptPath)
	if err != nil {
		lo.Error("failed to load prompt template: %s\n", err)
		return err
	}

	// check if the base branch exists
	ghservice := agithub.NewGitHubService(conf.Agent.GitHub.Owner, workRepository, gh, lo)
	if _, err = ghservice.GetBranch(baseBranch); err != nil {
		return err
	}

	functions.InitializeFunctions(
		*conf.Agent.GitHub.NoSubmit,
		ghservice,
		conf.Agent.AllowFunctions,
	)
	lo.Info("allowed functions: %s\n", strings.Join(util.Map(
		functions.AllFunctions(),
		func(e functions.Function) string { return e.Name.String() },
	), ","))
	lo.Info("agents make a pull request to %s/%s\n", conf.Agent.GitHub.Owner, workRepository)

	submitServiceCaller := agithub.NewSubmitFileGitHubService(conf.Agent.GitHub.Owner, workRepository, gh, lo).
		Caller(ctx, functions.SubmitFilesServiceInput{
			BaseBranch: baseBranch,
			GitEmail:   conf.Agent.Git.UserEmail,
			GitName:    conf.Agent.Git.UserName,
			PRLabels:   conf.Agent.GitHub.PRLabels,
		})

	dataStore := store.NewStore()

	parameter := Parameter{
		MaxSteps: conf.Agent.MaxSteps,
		Model:    conf.Agent.Model,
	}

	prompt, err := libprompt.BuildRequirementPrompt(promptTemplate, conf.Language, issue)
	if err != nil {
		lo.Error("failed build requirement prompt: %s\n", err)
		return err
	}
	requirementAgent, err := RunAgent("requirementAgent",
		prompt, submitServiceCaller, parameter, lo, &dataStore, llmForwarder, PlanTools())
	if err != nil {
		lo.Error("requirement agent failed: %s\n", err)
		return err
	}

	instruction := requirementAgent.LastHistory().RawContent
	prompt, err = libprompt.BuildDeveloperPrompt(promptTemplate, conf.Language, loaderr, issue.Path, instruction)
	if err != nil {
		lo.Error("failed build developer prompt: %s\n", err)
		return err
	}
	developerAgent, err := RunAgent("developerAgent",
		prompt, submitServiceCaller, parameter, lo, &dataStore, llmForwarder, functions.AllFunctions())
	if err != nil {
		lo.Error("developer agent failed: %s\n", err)
		return err
	}

	if conf.Agent.ReviewAgents == 0 {
		lo.Info("skip review agents\n")
		lo.Info("agents finished work\n")
		return nil
	}

	if s := dataStore.GetSubmission(store.LastSubmissionKey); s == nil {
		lo.Error("submission is not found\n")
		return err
	}
	submittedPRNumber := dataStore.GetSubmission(store.LastSubmissionKey).PullRequestNumber

	prompt, err = libprompt.BuildReviewManagerPrompt(promptTemplate, conf, issue, util.Map(developerAgent.ChangedFiles(), func(f store.File) string { return f.Path }))
	if err != nil {
		lo.Error("failed to build review manager prompt: %s\n", err)
		return err
	}
	reviewManager, err := RunAgent(
		"reviewManagerAgent",
		prompt,
		submitServiceCaller,
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
		prpt, err := libprompt.BuildReviewerPrompt(promptTemplate, conf.Language, submittedPRNumber, p.Prompt)
		if err != nil {
			lo.Error("failed to build reviewer prompt: %s\n", err)
			return err
		}

		reviewer, err := RunAgent(
			p.AgentName,
			prpt,
			submitServiceCaller,
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

func RunAgent(
	name string,
	prompt libprompt.Prompt,
	submitServiceCaller functions.SubmitFilesCallerType,
	parameter Parameter,
	lo logger.Logger,
	dataStore *store.Store,
	llmForwarder models.LLMForwarder,
	tools []functions.Function,
) (Agent, error) {
	ag := NewAgent(
		parameter,
		name,
		lo,
		submitServiceCaller,
		prompt,
		llmForwarder,
		dataStore,
		tools,
	)

	if _, err := ag.Work(); err != nil {
		lo.Error("requirement agent failed: %s\n", err)
		return Agent{}, err
	}

	return ag, nil
}
