package models_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/clover0/issue-agent/core"
	"github.com/clover0/issue-agent/models"
	"github.com/clover0/issue-agent/test/assert"
	"github.com/clover0/issue-agent/test/loggertest"
)

func TestSelectForwarder(t *testing.T) {
	t.Parallel()

	// TODO: Make it clear which environment variables need to be set.
	_ = os.Setenv("ANTHROPIC_API_KEY", "test")
	_ = os.Setenv("OPENAI_API_KEY", "test")

	mockLogger := loggertest.NewTestLogger()

	tests := map[string]struct {
		model    string
		wantErr  bool
		wantType core.LLMForwarder
	}{
		"AWS Bedrock model": {
			model:    "anthropic.claude-3-5-sonnet-v1",
			wantErr:  false,
			wantType: models.BedrockLLMForwarder{},
		},
		"OpenAI model": {
			model:    "gpt-4",
			wantErr:  false,
			wantType: models.OpenAILLMForwarder{},
		},
		"Anthropic model": {
			model:    "claude-3",
			wantErr:  false,
			wantType: models.AnthropicLLMForwarder{},
		},
		"Empty model": {
			model:    "",
			wantErr:  true,
			wantType: nil,
		},
		"Unsupported model": {
			model:    "unknown-model",
			wantErr:  true,
			wantType: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			forwarder, err := models.SelectForwarder(mockLogger, tt.model)

			if tt.wantErr {
				assert.HasError(t, err)
				assert.Nil(t, forwarder)
				return
			}

			assert.NoError(t, err)
			gotType := reflect.TypeOf(forwarder)
			wantType := reflect.TypeOf(tt.wantType)
			if gotType != wantType {
				t.Errorf("got: %v, want: %v", gotType, wantType)
			}
		})
	}
}
