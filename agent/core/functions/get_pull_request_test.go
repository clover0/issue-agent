package functions_test

import (
	"testing"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/test/assert"
)

func TestGetPullRequestOutput_ToLLMString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    functions.GetPullRequestOutput
		expected string
	}{
		"valid pull request information": {
			input: functions.GetPullRequestOutput{
				PRNumber: "123",
				Head:     "feature-branch",
				Base:     "main",
				RawDiff:  "diff --git a/file.txt b/file.txt\n...",
				Title:    "Feature: Add new functionality",
				Content:  "This is an amazing feature",
			},
			expected: `
<pr-number>
123
</pr-number>

<pull-request-title>
Feature: Add new functionality
</pull-request-title>

<pull-request-description>
This is an amazing feature
</pull-request-description>

<pull-request-diff>
diff --git a/file.txt b/file.txt
...
</pull-request-diff>
`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := tt.input.ToLLMString()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetReviewOutput_ToLLMString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    functions.GetReviewOutput
		expected string
	}{
		"valid review information": {
			input: functions.GetReviewOutput{
				IssuesNumber: "123",
				Path:         "src/main.go",
				StartLine:    10,
				EndLine:      15,
				Content:      "This is content",
			},
			expected: `
The following file information received a code review.

# Review information
* Review file path: src/main.go
* Review start line number: 10
* Review end line number: 15

# Review content
This is content
`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := tt.input.ToLLMString()
			assert.Equal(t, tt.expected, result)
		})
	}
}
