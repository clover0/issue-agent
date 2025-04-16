package models

import "strings"

const (
	Claude3_7MaxTokens = 64000
	Claude3_5MaxTokens = 8192
)

func ClaudeMaxTokens(model string) int {
	// https://docs.anthropic.com/en/docs/about-claude/models/all-models#model-comparison-table
	if strings.Contains(model, "claude-3-7-sonnet-") {
		return Claude3_7MaxTokens
	}

	// Claude-3-5.
	// Other models are not supported.
	return Claude3_5MaxTokens
}
