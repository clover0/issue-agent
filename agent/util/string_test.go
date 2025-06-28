package util

import (
	"testing"

	"github.com/clover0/issue-agent/test/assert"
)

func TestTruncateLines(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		text      string
		keepStart int
		keepEnd   int
		want      string
	}{
		"empty string": {
			text:      "",
			keepStart: 2,
			keepEnd:   2,
			want:      "",
		},
		"single line": {
			text:      "single line",
			keepStart: 2,
			keepEnd:   2,
			want:      "single line",
		},
		"few lines, no truncation needed": {
			text:      "line1\nline2\nline3\nline4",
			keepStart: 2,
			keepEnd:   2,
			want:      "line1\nline2\nline3\nline4",
		},
		"many lines, truncation needed": {
			text:      "line1\nline2\nline3\nline4\nline5\nline6",
			keepStart: 2,
			keepEnd:   2,
			want:      "line1\nline2\n...\nline5\nline6",
		},
		"negative keepStart": {
			text:      "line1\nline2\nline3\nline4\nline5",
			keepStart: -1,
			keepEnd:   2,
			want:      "line1\nline2\nline3\nline4\nline5",
		},
		"negative keepEnd": {
			text:      "line1\nline2\nline3\nline4\nline5",
			keepStart: 2,
			keepEnd:   -1,
			want:      "line1\nline2\nline3\nline4\nline5",
		},
		"zero keepStart, positive keepEnd": {
			text:      "line1\nline2\nline3\nline4\nline5",
			keepStart: 0,
			keepEnd:   2,
			want:      "...\nline4\nline5",
		},
		"positive keepStart, zero keepEnd": {
			text:      "line1\nline2\nline3\nline4\nline5",
			keepStart: 2,
			keepEnd:   0,
			want:      "line1\nline2\n...",
		},
		"zero keepStart, zero keepEnd": {
			text:      "line1\nline2\nline3\nline4\nline5",
			keepStart: 0,
			keepEnd:   0,
			want:      "...",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := TruncateLines(tt.text, tt.keepStart, tt.keepEnd, "...")

			assert.Equal(t, tt.want, result)
		})
	}
}
