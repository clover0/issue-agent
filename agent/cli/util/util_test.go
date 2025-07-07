package util

import (
	"testing"

	"github.com/clover0/issue-agent/test/assert"
)

func TestParseArgFlags(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input     []string
		wantArg   string
		wantFlags []string
	}{
		"empty input": {
			input:     []string{},
			wantArg:   "",
			wantFlags: []string{},
		},
		"only argument": {
			input:     []string{"arg1"},
			wantArg:   "arg1",
			wantFlags: []string{},
		},
		"argument with one flag": {
			input:     []string{"arg1", "--flag1"},
			wantArg:   "arg1",
			wantFlags: []string{"--flag1"},
		},
		"argument with multiple flags": {
			input:     []string{"arg1", "--flag1", "value1", "--flag3"},
			wantArg:   "arg1",
			wantFlags: []string{"--flag1", "value1", "--flag3"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotArg, gotFlags := ParseArgFlags(tt.input)

			assert.Equal(t, gotArg, tt.wantArg)
			assert.EqualStringSlices(t, gotFlags, tt.wantFlags)

		})
	}
}
