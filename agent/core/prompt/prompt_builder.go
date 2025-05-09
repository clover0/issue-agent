package prompt

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/core/functions"
)

type Prompt struct {
	SystemPrompt    string
	StartUserPrompt string
}

func BuildRequirementPrompt(promptTpl PromptTemplate, language string, baseBranch string, issue functions.GetIssueOutput) (Prompt, error) {
	return BuildPrompt(promptTpl, "planner", map[string]any{
		"language":     language,
		"issueTitle":   issue.Title,
		"issueContent": issue.Content,
		"issueNumber":  issue.Path,
		"baseBranch":   baseBranch,
	})
}

func BuildDeveloperPrompt(promptTpl PromptTemplate, language string, baseBranch string, issue functions.GetIssueOutput, instruction string) (Prompt, error) {
	return BuildPrompt(promptTpl, "developer", map[string]any{
		"language":     language,
		"issueTitle":   issue.Title,
		"issueContent": issue.Content,
		"issueNumber":  issue.Path,
		"instruction":  instruction,
		"baseBranch":   baseBranch,
	})
}

func BuildCommentReactorPrompt(promptTpl PromptTemplate, language string,
	comment functions.GetCommentOutput,
	pr functions.GetPullRequestOutput) (Prompt, error) {
	return BuildPrompt(promptTpl, "comment-reactor", map[string]any{
		"language":      language,
		"workingBranch": pr.Head,
		"prNumber":      pr.PRNumber,
		"issueNumber":   comment.IssueNumber,
		"comment":       comment.Content,
		"prLLMString":   pr.ToLLMString(),
	})
}

func BuildReviewManagerPrompt(promptTpl PromptTemplate, cnf config.Config, issue functions.GetIssueOutput, changedFilesPath []string, baseBranch string) (Prompt, error) {
	m := make(map[string]any)

	m["language"] = cnf.Language
	m["filePaths"] = changedFilesPath
	m["issue"] = issue.Content
	m["reviewAgents"] = cnf.Agent.ReviewAgents
	m["baseBranch"] = baseBranch

	m["noFiles"] = ""
	if len(changedFilesPath) == 0 {
		m["noFiles"] = "no changed files"
	}

	return BuildPrompt(promptTpl, "review-manager", m)
}

func BuildReviewerPrompt(promptTpl PromptTemplate, language string, prNumber int, reviewerPrompt string) (Prompt, error) {
	return BuildPrompt(promptTpl, "reviewer", map[string]any{
		"language":       language,
		"prNumber":       prNumber,
		"reviewerPrompt": reviewerPrompt,
	})
}

func BuildPrompt(promptTpl PromptTemplate, templateName string, templateMap map[string]any) (Prompt, error) {
	var prpt Prompt
	for _, p := range promptTpl.Agents {
		if p.Name == templateName {
			prpt = Prompt{
				SystemPrompt:    p.SystemTemplate,
				StartUserPrompt: p.UserTemplate,
			}
			break
		}
	}
	if prpt.StartUserPrompt == "" {
		return Prompt{}, fmt.Errorf("failed to find %s prompt. you must have  name=%s prompt in the prompt template", templateName, templateName)
	}

	systemPrompt, err := parseTemplate(prpt.SystemPrompt, templateMap)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to parse system prompt: %w", err)
	}

	userPrompt, err := parseTemplate(prpt.StartUserPrompt, templateMap)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to parse user prompt: %w", err)
	}

	return Prompt{
		SystemPrompt:    systemPrompt,
		StartUserPrompt: userPrompt,
	}, nil
}

func parseTemplate(templateStr string, values map[string]any) (string, error) {
	tpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", err
	}

	tplbuff := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(tplbuff, values); err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return tplbuff.String(), nil
}
