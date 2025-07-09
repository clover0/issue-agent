package prompt_test

import (
	"testing"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/core/prompt"
	"github.com/clover0/issue-agent/test/assert"
)

func TestBuildPlanningPromptPrompt(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.Template{
			Agents: []prompt.AgentPromptTemplate{
				{
					Name:           "planner",
					SystemTemplate: "System: {{.language}} issue {{.issueNumber}}",
					UserTemplate:   "User: {{.issueTitle}} - {{.issueContent}} - {{.baseBranch}}",
				},
			},
		}

		issue := functions.GetIssueOutput{
			Path:    "123",
			Title:   "Test Issue",
			Content: "This is a test issue",
		}

		p, err := prompt.BuildPlanningPrompt(promptTpl, "English", "main", issue)

		assert.Nil(t, err)
		assert.Equal(t, p.SystemPrompt, "System: English issue 123")
		assert.Equal(t, p.StartUserPrompt, "User: Test Issue - This is a test issue - main")
	})

	t.Run("template not found", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.Template{
			Agents: []prompt.AgentPromptTemplate{
				{
					Name:           "other-agent",
					SystemTemplate: "other system template",
					UserTemplate:   "other user template",
				},
			},
		}

		issue := functions.GetIssueOutput{
			Path:    "123",
			Title:   "Test Issue",
			Content: "This is a test issue",
		}

		_, err := prompt.BuildPlanningPrompt(promptTpl, "English", "main", issue)

		assert.HasError(t, err)
	})
}

func TestBuildDeveloperPrompt(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.Template{
			Agents: []prompt.AgentPromptTemplate{
				{
					Name:           "developer",
					SystemTemplate: "System: {{.language}} issue {{.issueNumber}} with instruction",
					UserTemplate:   "User: {{.issueTitle}} - {{.issueContent}} - {{.instruction}} - {{.baseBranch}}",
				},
			},
		}

		issue := functions.GetIssueOutput{
			Path:    "123",
			Title:   "Test Issue",
			Content: "This is a test issue",
		}

		instruction := "Implement the feature"

		p, err := prompt.BuildDeveloperPrompt(promptTpl, "English", "main", issue, instruction)

		assert.Nil(t, err)
		assert.Equal(t, p.SystemPrompt, "System: English issue 123 with instruction")
		assert.Equal(t, p.StartUserPrompt, "User: Test Issue - This is a test issue - Implement the feature - main")
	})

	t.Run("template not found", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.Template{
			Agents: []prompt.AgentPromptTemplate{
				{
					Name:           "other-agent",
					SystemTemplate: "other system template",
					UserTemplate:   "other user template",
				},
			},
		}

		issue := functions.GetIssueOutput{
			Path:    "123",
			Title:   "Test Issue",
			Content: "This is a test issue",
		}

		instruction := "Implement the feature"

		_, err := prompt.BuildDeveloperPrompt(promptTpl, "English", "main", issue, instruction)

		assert.HasError(t, err)
	})
}

func TestFindPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.Template{
			Agents: []prompt.AgentPromptTemplate{
				{
					Name:           "test-template-first",
					SystemTemplate: "test system template first",
					UserTemplate:   "test user template first",
				},
				{
					Name:           "test-template-second",
					SystemTemplate: "test system template second",
					UserTemplate:   "test user template second",
				},
			},
		}

		p, err := prompt.FindPromptTemplate(promptTpl, "test-template-second")

		assert.Nil(t, err)
		assert.Equal(t, p.SystemTemplate, "test system template second")
		assert.Equal(t, p.UserTemplate, "test user template second")
	})

	t.Run("success on find one", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.Template{
			Agents: []prompt.AgentPromptTemplate{
				{
					Name:           "test-template",
					SystemTemplate: "test system template first",
					UserTemplate:   "test user template first",
				},
				{
					Name:           "test-template",
					SystemTemplate: "test system template second",
					UserTemplate:   "test user template second",
				},
			},
		}

		p, err := prompt.FindPromptTemplate(promptTpl, "test-template")

		assert.Nil(t, err)
		assert.Equal(t, p.SystemTemplate, "test system template first")
		assert.Equal(t, p.UserTemplate, "test user template first")
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.Template{
			Agents: []prompt.AgentPromptTemplate{
				{
					Name:           "other-agent",
					SystemTemplate: "other system template",
					UserTemplate:   "other user template",
				},
			},
		}

		_, err := prompt.FindPromptTemplate(promptTpl, "non-existent-agent")

		assert.HasError(t, err)
	})
}

func TestBuildPrompt(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.AgentPromptTemplate{
			Name:           "test-template",
			SystemTemplate: "Hello {{.name}}",
			UserTemplate:   "How are you {{.name}}?",
		}

		templateMap := map[string]any{
			"name": "World",
		}

		p, err := prompt.BuildPrompt(promptTpl, templateMap)

		assert.Nil(t, err)
		assert.Equal(t, p.SystemPrompt, "Hello World")
		assert.Equal(t, p.StartUserPrompt, "How are you World?")
	})

	t.Run("error in system template", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.AgentPromptTemplate{
			Name:           "test-template",
			SystemTemplate: "Hello {{.invalid}",
			UserTemplate:   "How are you {{.name}}?",
		}

		templateMap := map[string]any{
			"name": "World",
		}

		_, err := prompt.BuildPrompt(promptTpl, templateMap)

		assert.HasError(t, err)
	})

	t.Run("error in user template", func(t *testing.T) {
		t.Parallel()

		promptTpl := prompt.AgentPromptTemplate{
			Name:           "test-template",
			SystemTemplate: "Hello {{.name}}",
			UserTemplate:   "How are you {{.invalid}",
		}

		templateMap := map[string]any{
			"name": "World",
		}

		_, err := prompt.BuildPrompt(promptTpl, templateMap)

		assert.HasError(t, err)
	})
}
