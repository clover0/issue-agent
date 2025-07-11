package prompt_test

import (
	"testing"

	"github.com/clover0/issue-agent/core/prompt"
	"github.com/clover0/issue-agent/test/assert"
)

func TestParseTemplate(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Name    string
		Value   int
		Enabled bool
	}

	tests := map[string]struct {
		templateStr string
		values      TestStruct
		want        string
		wantErr     bool
	}{
		"simple template": {
			templateStr: "Hello, {{.Name}}!",
			values: TestStruct{
				Name: "World",
			},
			want:    "Hello, World!",
			wantErr: false,
		},
		"multiple values": {
			templateStr: "Name: {{.Name}}, Value: {{.Value}}, Enabled: {{.Enabled}}",
			values: TestStruct{
				Name:    "Test",
				Value:   42,
				Enabled: true,
			},
			want:    "Name: Test, Value: 42, Enabled: true",
			wantErr: false,
		},
		"with newlines": {
			templateStr: `
Hello, {{.Name}}!
Your value is {{.Value}}.
Enabled: {{.Enabled}}
`,
			values: TestStruct{
				Name:    "User",
				Value:   100,
				Enabled: false,
			},
			want: `
Hello, User!
Your value is 100.
Enabled: false
`,
			wantErr: false,
		},
		"invalid template": {
			templateStr: "Hello, {{.InvalidField}}!",
			values: TestStruct{
				Name: "World",
			},
			want:    "",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := prompt.ParseTemplate(tt.templateStr, tt.values)

			if tt.wantErr {
				assert.HasError(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
