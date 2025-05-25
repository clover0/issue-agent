package models_test

import (
	"testing"

	"github.com/clover0/issue-agent/models"
	"github.com/clover0/issue-agent/test/assert"
)

func TestClaudeMaxTokens(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		model string
		want  int
	}{
		"AWS Bedrock: anthropic.claude-3-7-sonnet-20250219-v1:0": {
			model: "anthropic.claude-3-7-sonnet-20250219-v1:0",
			want:  models.ClaudeMaxTokens64K,
		},
		"Anthropic: claude-3-5 model": {
			model: "claude-3-5-sonnet-20240620",
			want:  models.ClaudeDefaultMaxTokens,
		},
		"unsupported model": {
			model: "unknown-model",
			want:  models.ClaudeDefaultMaxTokens,
		},
		"empty model string": {
			model: "",
			want:  models.ClaudeDefaultMaxTokens,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := models.ClaudeMaxOutputTokens(tt.model)

			assert.Equal(t, tt.want, result)
		})
	}
}
