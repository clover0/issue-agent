package createpr_test

import (
	"testing"

	"github.com/clover0/issue-agent/cli/command/common"
	"github.com/clover0/issue-agent/cli/command/createpr"
	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/test/assert"
)

func TestMergeGitHubArg(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input *createpr.CreatePRInput
		arg   createpr.ArgGitHubCreatePR
		want  *createpr.CreatePRInput
	}{
		"merge valid ArgGitHubCreatePR": {
			input: &createpr.CreatePRInput{
				GitHubOwner:       "",
				WorkRepository:    "",
				GithubIssueNumber: "",
			},
			arg: createpr.ArgGitHubCreatePR{
				Owner:       "newOwner",
				Repository:  "newRepo",
				IssueNumber: "456",
			},
			want: &createpr.CreatePRInput{
				GitHubOwner:       "newOwner",
				WorkRepository:    "newRepo",
				GithubIssueNumber: "456",
			},
		},
		"merge with existing values": {
			input: &createpr.CreatePRInput{
				GitHubOwner:       "existingOwner",
				WorkRepository:    "existingRepo",
				GithubIssueNumber: "123",
			},
			arg: createpr.ArgGitHubCreatePR{
				Owner:       "newOwner",
				Repository:  "newRepo",
				IssueNumber: "456",
			},
			want: &createpr.CreatePRInput{
				GitHubOwner:       "newOwner",
				WorkRepository:    "newRepo",
				GithubIssueNumber: "456",
			},
		},
		"merge with empty ArgGitHubCreatePR": {
			input: &createpr.CreatePRInput{
				GitHubOwner:       "existingOwner",
				WorkRepository:    "existingRepo",
				GithubIssueNumber: "123",
			},
			arg: createpr.ArgGitHubCreatePR{
				Owner:       "",
				Repository:  "",
				IssueNumber: "",
			},
			want: &createpr.CreatePRInput{
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
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMergeConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input  *createpr.CreatePRInput
		config config.Config
		want   config.Config
	}{
		"merge with empty input": {
			input: &createpr.CreatePRInput{
				Common:            &common.CommonInput{},
				GitHubOwner:       "",
				WorkRepository:    "",
				GithubIssueNumber: "",
				BaseBranch:        "",
			},
			config: config.Config{
				LogLevel: "info",
				Language: "English",
				Agent: config.AgentConfig{
					Model: "gpt-4",
					GitHub: config.GitHubConfig{
						Owner: "original-owner",
					},
				},
			},
			want: config.Config{
				LogLevel: "info",
				Language: "English",
				Agent: config.AgentConfig{
					Model: "gpt-4",
					GitHub: config.GitHubConfig{
						Owner: "original-owner",
					},
				},
			},
		},
		"merge with all fields populated": {
			input: &createpr.CreatePRInput{
				Common: &common.CommonInput{
					LogLevel: "debug",
					Language: "Japanese",
					Model:    "claude",
				},
				GitHubOwner:       "new-owner",
				WorkRepository:    "repo",
				GithubIssueNumber: "123",
				BaseBranch:        "main",
				Reviewers:         []string{"reviewer1", "reviewer2"},
				TeamReviewers:     []string{"team1", "team2"},
			},
			config: config.Config{
				LogLevel: "info",
				Language: "English",
				Agent: config.AgentConfig{
					Model: "gpt-4",
					GitHub: config.GitHubConfig{
						Owner: "original-owner",
					},
				},
			},
			want: config.Config{
				LogLevel: "debug",
				Language: "Japanese",
				Agent: config.AgentConfig{
					Model: "claude",
					GitHub: config.GitHubConfig{
						Owner:         "new-owner",
						Reviewers:     []string{"reviewer1", "reviewer2"},
						TeamReviewers: []string{"team1", "team2"},
					},
				},
			},
		},
		"merge with some fields populated": {
			input: &createpr.CreatePRInput{
				Common: &common.CommonInput{
					LogLevel: "debug",
					Model:    "",
				},
				GitHubOwner:       "new-owner",
				WorkRepository:    "repo",
				GithubIssueNumber: "123",
				BaseBranch:        "main",
			},
			config: config.Config{
				LogLevel: "info",
				Language: "English",
				Agent: config.AgentConfig{
					Model: "gpt-4",
					GitHub: config.GitHubConfig{
						Owner:         "original-owner",
						TeamReviewers: []string{"original-team"},
					},
				},
			},
			want: config.Config{
				LogLevel: "debug",
				Language: "English",
				Agent: config.AgentConfig{
					Model: "gpt-4",
					GitHub: config.GitHubConfig{
						Owner:         "new-owner",
						TeamReviewers: []string{"original-team"},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := tt.input.MergeConfig(tt.config)
			assert.Equal(t, result, tt.want)
		})
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   createpr.CreatePRInput
		wantErr bool
	}{
		"valid input": {
			input: createpr.CreatePRInput{
				GitHubOwner:       "owner",
				WorkRepository:    "repo",
				BaseBranch:        "main",
				GithubIssueNumber: "123",
			},
			wantErr: false,
		},
		"missing required fields": {
			input: createpr.CreatePRInput{
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
		want    createpr.ArgGitHubCreatePR
		wantErr bool
	}{
		"valid input": {
			input: "owner/repo/issues/123",
			want: createpr.ArgGitHubCreatePR{
				Owner:       "owner",
				Repository:  "repo",
				IssueNumber: "123",
			},
			wantErr: false,
		},
		"invalid input: missing `issues` segment": {
			input:   "owner/repo/123",
			want:    createpr.ArgGitHubCreatePR{},
			wantErr: true,
		},
		"invalid input: too many segments": {
			input:   "owner/repo/issues/123/extra",
			want:    createpr.ArgGitHubCreatePR{},
			wantErr: true,
		},
		"invalid input: not enough segments (missing owner)": {
			input:   "repo/issues/123",
			want:    createpr.ArgGitHubCreatePR{},
			wantErr: true,
		},
		"invalid input: not enough segments (missing repository)": {
			input:   "owner/issues/123",
			want:    createpr.ArgGitHubCreatePR{},
			wantErr: true,
		},
		"invalid input: empty string": {
			input:   "",
			want:    createpr.ArgGitHubCreatePR{},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := createpr.ParseCreatePRGitHubArg(tt.input)

			if tt.wantErr {
				assert.HasError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
