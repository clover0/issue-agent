package models

import "strings"

const (
	ClaudeMaxTokens64K     = 1024 * 64
	ClaudeOpusMaxTokens32K = 1024 * 32
	ClaudeDefaultMaxTokens = 1024 * 8
)

func ClaudeMaxOutputTokens(model string) int {
	// https://docs.anthropic.com/en/docs/about-claude/models/all-models#model-comparison-table
	if strings.Contains(model, "claude-3-7-sonnet-") {
		return ClaudeMaxTokens64K
	}
	if strings.Contains(model, "claude-sonnet-4") {
		return ClaudeMaxTokens64K
	}
	if strings.Contains(model, "claude-opus-4") {
		return ClaudeOpusMaxTokens32K
	}

	// Claude-3-5.
	// Other models are not supported.
	return ClaudeDefaultMaxTokens
}
