package react_test

import (
	"testing"

	"github.com/clover0/issue-agent/cli/command/react"
	"github.com/clover0/issue-agent/test/assert"
)

func TestBindReactGitHubArg(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    react.ArgGitHubReact
		wantErr bool
	}{
		"valid: pull request comment input": {
			input: "owner/repo/issues/comments/123456",
			want: react.ArgGitHubReact{
				ReactType:  react.Comment,
				Owner:      "owner",
				Repository: "repo",
				PRNumber:   "",
				CommentID:  "123456",
				ReviewID:   "",
			},
			wantErr: false,
		},
		"valid: pull request review comment input": {
			input: "owner/repo/pulls/comments/789012",
			want: react.ArgGitHubReact{
				ReactType:  react.ReviewComment,
				Owner:      "owner",
				Repository: "repo",
				PRNumber:   "",
				CommentID:  "",
				ReviewID:   "789012",
			},
			wantErr: false,
		},
		"invalid input: wrong format": {
			input:   "owner/repo/something/comments/123456",
			want:    react.ArgGitHubReact{},
			wantErr: true,
		},
		"invalid input: missing comment ID": {
			input:   "owner/repo/issues/comments/",
			want:    react.ArgGitHubReact{},
			wantErr: true,
		},
		"invalid input: too many segments": {
			input:   "owner/repo/issues/comments/123456/extra",
			want:    react.ArgGitHubReact{},
			wantErr: true,
		},
		"invalid input: not enough segments": {
			input:   "owner/issues/comments/123456",
			want:    react.ArgGitHubReact{},
			wantErr: true,
		},
		"invalid input: empty string": {
			input:   "",
			want:    react.ArgGitHubReact{},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := react.BindReactGitHubArg(tt.input)

			if tt.wantErr {
				assert.HasError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
