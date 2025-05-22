package prompt

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/core/functions"
)

func BuildRequirementPrompt(promptTpl Template, language string, baseBranch string, issue functions.GetIssueOutput) (Prompt, error) {
	tmpl, err := FindPromptTemplate(promptTpl, "requirement")
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to find requirement prompt template: %w", err)
	}

	return BuildPrompt(tmpl, map[string]any{
		"language":     language,
		"issueTitle":   issue.Title,
		"issueContent": issue.Content,
		"issueNumber":  issue.Path,
		"baseBranch":   baseBranch,
	})
}

func BuildDeveloperPrompt(promptTpl Template, language string, baseBranch string, issue functions.GetIssueOutput, instruction string) (Prompt, error) {
	tmpl, err := FindPromptTemplate(promptTpl, "developer")
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to find developer prompt template: %w", err)
	}

	return BuildPrompt(tmpl, map[string]any{
		"language":     language,
		"issueTitle":   issue.Title,
		"issueContent": issue.Content,
		"issueNumber":  issue.Path,
		"instruction":  instruction,
		"baseBranch":   baseBranch,
	})
}

func BuildCommentReactorPrompt(promptTpl Template, language string,
	comment functions.GetCommentOutput,
	pr functions.GetPullRequestOutput) (Prompt, error) {
	tmpl, err := FindPromptTemplate(promptTpl, "comment-reactor")
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to find comment-reactor prompt template: %w", err)
	}

	return BuildPrompt(tmpl, map[string]any{
		"language":      language,
		"workingBranch": pr.Head,
		"prNumber":      pr.PRNumber,
		"issueNumber":   comment.IssueNumber,
		"comment":       comment.Content,
		"prLLMString":   pr.ToLLMString(),
	})
}

func BuildReviewManagerPrompt(promptTpl Template, cnf config.Config, issue functions.GetIssueOutput, changedFilesPath []string, baseBranch string) (Prompt, error) {
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

	tmpl, err := FindPromptTemplate(promptTpl, "review-manager")
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to find review-manager prompt template: %w", err)
	}

	return BuildPrompt(tmpl, m)
}

func BuildReviewerPrompt(promptTpl Template, language string, prNumber int, reviewerPrompt string) (Prompt, error) {
	tmpl, err := FindPromptTemplate(promptTpl, "reviewer")
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to find reviewer prompt template: %w", err)
	}

	return BuildPrompt(tmpl, map[string]any{
		"language":       language,
		"prNumber":       prNumber,
		"reviewerPrompt": reviewerPrompt,
	})
}

func FindPromptTemplate(promptTpl Template, name string) (AgentPromptTemplate, error) {
	for _, p := range promptTpl.Agents {
		if p.Name == name {
			return p, nil
		}
	}

	return AgentPromptTemplate{}, fmt.Errorf("failed to find %s prompt. you must have  name=%s prompt in the prompt template", name, name)
}

func BuildPrompt(promptTpl AgentPromptTemplate, templateMap map[string]any) (Prompt, error) {
	systemPrompt, err := parseTemplate(promptTpl.SystemTemplate, templateMap)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to parse system prompt: %w", err)
	}

	userPrompt, err := parseTemplate(promptTpl.UserTemplate, templateMap)
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
