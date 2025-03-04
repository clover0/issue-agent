package models

import (
	"context"
	"fmt"
	"os"

	"github.com/clover0/issue-agent/core"
	"github.com/clover0/issue-agent/functions"
	"github.com/clover0/issue-agent/logger"
)

type OpenAILLMForwarder struct {
	openai OpenAI
}

func NewOpenAILLMForwarder(l logger.Logger) (core.LLMForwarder, error) {
	apiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	return OpenAILLMForwarder{
		openai: NewOpenAI(l, apiKey),
	}, nil
}

func (o OpenAILLMForwarder) StartForward(input core.StartCompletionInput) ([]core.LLMMessage, error) {
	return o.openai.StartCompletion(
		context.TODO(),
		core.StartCompletionInput{
			Model:           input.Model,
			SystemPrompt:    input.SystemPrompt,
			StartUserPrompt: input.StartUserPrompt,
			Functions:       functions.AllFunctions(),
		},
	)
}

func (o OpenAILLMForwarder) ForwardLLM(
	ctx context.Context,
	input core.StartCompletionInput,
	llmContexts []core.ReturnToLLMContext,
	history []core.LLMMessage,
) ([]core.LLMMessage, error) {
	return o.openai.ContinueCompletion(ctx, input, llmContexts, history)
}

func (o OpenAILLMForwarder) ForwardStep(ctx context.Context, history []core.LLMMessage) core.Step {
	return o.openai.CompletionNextStep(ctx, history)
}
