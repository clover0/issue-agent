package models

import (
	"fmt"
	"strings"

	"github.com/clover0/issue-agent/core"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/util"
)

func SelectForwarder(lo logger.Logger, model string) (core.LLMForwarder, error) {
	if util.IsAWSBedrockModel(model) {
		return NewBedrockLLMForwarder(lo)
	}
	if strings.HasPrefix(model, "gpt") {
		return NewOpenAILLMForwarder(lo)
	}

	if strings.HasPrefix(model, "claude") {
		return NewAnthropicLLMForwarder(lo)
	}

	if model == "" {
		return nil, fmt.Errorf("model is not specified")
	}

	return nil, fmt.Errorf("SelectForwarder: model %s is not supported", model)
}
