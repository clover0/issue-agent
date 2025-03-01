package prompt

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/loader"
)

type Prompt struct {
	SystemPrompt    string
	StartUserPrompt string
}

func BuildRequirementPrompt(promptTpl PromptTemplate, language string, baseBranch string, issue loader.Issue) (Prompt, error) {
	return BuildPrompt(promptTpl, "planner", map[string]any{
		"language":    language,
		"issue":       issue.Content,
		"issueNumber": issue.Path,
		"baseBranch":  baseBranch,
	})
}

func BuildDeveloperPrompt(promptTpl PromptTemplate, language string, baseBranch string, issueLoader loader.Loader, issueNumber string, instruction string) (Prompt, error) {
	// TODO: separate issueLoader and issue from this
	iss, err := issueLoader.LoadIssue(context.TODO(), issueNumber)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to load issue: %w", err)
	}

	return BuildPrompt(promptTpl, "developer", map[string]any{
		"language":    language,
		"issue":       iss.Content,
		"issueNumber": issueNumber,
		"instruction": instruction,
		"baseBranch":  baseBranch,
	})
}

func BuildReviewManagerPrompt(promptTpl PromptTemplate, cnf config.Config, issue loader.Issue, changedFilesPath []string, baseBranch string) (Prompt, error) {
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
