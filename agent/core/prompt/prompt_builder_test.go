package prompt

import (
	"testing"

	"github.com/clover0/issue-agent/test/assert"
)

func TestFindPromptTemplate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		promptTpl := PromptTemplate{
			Agents: []struct {
				Name           string `yaml:"name"`
				SystemTemplate string `yaml:"system_prompt"`
				UserTemplate   string `yaml:"user_prompt"`
			}{
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

		prompt, err := FindPromptTemplate(promptTpl, "test-template-second")

		assert.Nil(t, err)
		assert.Equal(t, prompt.SystemPrompt, "test system template second")
		assert.Equal(t, prompt.StartUserPrompt, "test user template second")
	})

	t.Run("success on find one", func(t *testing.T) {
		promptTpl := PromptTemplate{
			Agents: []struct {
				Name           string `yaml:"name"`
				SystemTemplate string `yaml:"system_prompt"`
				UserTemplate   string `yaml:"user_prompt"`
			}{
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

		prompt, err := FindPromptTemplate(promptTpl, "test-template")

		assert.Nil(t, err)
		assert.Equal(t, prompt.SystemPrompt, "test system template first")
		assert.Equal(t, prompt.StartUserPrompt, "test user template first")
	})

	t.Run("not_found", func(t *testing.T) {
		promptTpl := PromptTemplate{
			Agents: []struct {
				Name           string `yaml:"name"`
				SystemTemplate string `yaml:"system_prompt"`
				UserTemplate   string `yaml:"user_prompt"`
			}{
				{
					Name:           "other-agent",
					SystemTemplate: "other system template",
					UserTemplate:   "other user template",
				},
			},
		}

		_, err := FindPromptTemplate(promptTpl, "non-existent-agent")

		assert.HasError(t, err)
	})
}
