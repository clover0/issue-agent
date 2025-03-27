package createpr

import (
	"testing"

	"github.com/clover0/issue-agent/test/assert"
)

func TestMergeGitHubArg(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    *CreatePRInput
		arg      ArgGitHubCreatePR
		expected *CreatePRInput
	}{
		"merge valid ArgGitHubCreatePR": {
			input: &CreatePRInput{
				GitHubOwner:       "",
				WorkRepository:    "",
				GithubIssueNumber: "",
			},
			arg: ArgGitHubCreatePR{
				Owner:       "newOwner",
				Repository:  "newRepo",
				IssueNumber: "456",
			},
			expected: &CreatePRInput{
				GitHubOwner:       "newOwner",
				WorkRepository:    "newRepo",
				GithubIssueNumber: "456",
			},
		},
		"merge with existing values": {
			input: &CreatePRInput{
				GitHubOwner:       "existingOwner",
				WorkRepository:    "existingRepo",
				GithubIssueNumber: "123",
			},
			arg: ArgGitHubCreatePR{
				Owner:       "newOwner",
				Repository:  "newRepo",
				IssueNumber: "456",
			},
			expected: &CreatePRInput{
				GitHubOwner:       "newOwner",
				WorkRepository:    "newRepo",
				GithubIssueNumber: "456",
			},
		},
		"merge with empty ArgGitHubCreatePR": {
			input: &CreatePRInput{
				GitHubOwner:       "existingOwner",
				WorkRepository:    "existingRepo",
				GithubIssueNumber: "123",
			},
			arg: ArgGitHubCreatePR{
				Owner:       "",
				Repository:  "",
				IssueNumber: "",
			},
			expected: &CreatePRInput{
				GitHubOwner:       "",
				WorkRepository:    "",
				GithubIssueNumber: "",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := tt.input.MergeGitHubArg(tt.arg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMergeConfig(t *testing.T) {
	t.Parallel()
	t.Skipf("TODO")
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   CreatePRInput
		wantErr bool
	}{
		"valid input": {
			input: CreatePRInput{
				GitHubOwner:       "owner",
				WorkRepository:    "repo",
				BaseBranch:        "main",
				GithubIssueNumber: "123",
			},
			wantErr: false,
		},
		"missing required fields": {
			input: CreatePRInput{
				GitHubOwner: "",
				BaseBranch:  "main",
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tt.input.Validate()

			if tt.wantErr {
				assert.HasError(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseGitHubArg(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    ArgGitHubCreatePR
		wantErr bool
	}{
		"valid input": {
			input: "owner/repo/issues/123",
			want: ArgGitHubCreatePR{
				Owner:       "owner",
				Repository:  "repo",
				IssueNumber: "123",
			},
			wantErr: false,
		},
		"invalid input: missing `issues` segment": {
			input:   "owner/repo/123",
			want:    ArgGitHubCreatePR{},
			wantErr: true,
		},
		"invalid input: too many segments": {
			input:   "owner/repo/issues/123/extra",
			want:    ArgGitHubCreatePR{},
			wantErr: true,
		},
		"invalid input: not enough segments (missing owner)": {
			input:   "repo/issues/123",
			want:    ArgGitHubCreatePR{},
			wantErr: true,
		},
		"invalid input: not enough segments (missing repository)": {
			input:   "owner/issues/123",
			want:    ArgGitHubCreatePR{},
			wantErr: true,
		},
		"invalid input: empty string": {
			input:   "",
			want:    ArgGitHubCreatePR{},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseCreatePRGitHubArg(tt.input)

			if tt.wantErr {
				assert.HasError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
