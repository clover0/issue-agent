package react

import (
	"testing"

	"github.com/clover0/issue-agent/test/assert"
)

func TestBindReactGitHubArg(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    ArgGitHubReact
		wantErr bool
	}{
		"valid: pull request comment input": {
			input: "owner/repo/issues/comments/123456",
			want: ArgGitHubReact{
				ReactType:  Comment,
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
			want: ArgGitHubReact{
				ReactType:  ReviewComment,
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
			want:    ArgGitHubReact{},
			wantErr: true,
		},
		"invalid input: missing comment ID": {
			input:   "owner/repo/issues/comments/",
			want:    ArgGitHubReact{},
			wantErr: true,
		},
		"invalid input: too many segments": {
			input:   "owner/repo/issues/comments/123456/extra",
			want:    ArgGitHubReact{},
			wantErr: true,
		},
		"invalid input: not enough segments": {
			input:   "owner/issues/comments/123456",
			want:    ArgGitHubReact{},
			wantErr: true,
		},
		"invalid input: empty string": {
			input:   "",
			want:    ArgGitHubReact{},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := BindReactGitHubArg(tt.input)

			if tt.wantErr {
				assert.HasError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
