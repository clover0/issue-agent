package core

import (
	"context"
	"fmt"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/util"
)

type StartCompletionInput struct {
	Model           string
	SystemPrompt    string
	StartUserPrompt string
	Functions       []functions.Function
}

type LLMForwarder interface {
	StartForward(input StartCompletionInput) ([]LLMMessage, error)
	ForwardLLM(
		ctx context.Context,
		input StartCompletionInput,
		llmContexts []ReturnToLLMContext,
		history []LLMMessage,
	) ([]LLMMessage, error)
	ForwardStep(ctx context.Context, history []LLMMessage) Step
}

type SelectForwarder = func(lo logger.Logger, model string) (LLMForwarder, error)

type LLMMessage struct {
	Role         MessageRole
	RawContent   string
	FinishReason MessageFinishReason

	// user to llm
	RespondToolCall ToolCall

	// llm to user
	ReturnedToolCalls []ToolCall

	// returned raw message struct from LLM API
	RawMessageStruct any

	// Usage saves LLM usage information
	// Only the usage response from LLM response message,
	// so Usage is stored in Message with Role = LLMAssistant or LLMTool.
	Usage LLMUsage
}

func (l LLMMessage) ShowAssistantMessage(out logger.Logger) {
	out.Info(fmt.Sprintf("finish_reason: %s, cache create token: %d, cache read token: %d, "+
		"input token: %d, output token: %d, total input token: %d, total output token: %d\n",
		l.FinishReason, l.Usage.CacheCreateToken, l.Usage.CacheReadToken,
		l.Usage.InputToken, l.Usage.OutputToken, l.Usage.TotalInputToken(), l.Usage.TotalOutputToken()))

	out.Debug("message: \n")
	out.Debug(fmt.Sprintf("%s\n", l.RawContent))
	out.Debug("tools:\n")
	for _, t := range l.ReturnedToolCalls {
		out.Debug(fmt.Sprintf("id: %s, function_name:%s, function_args:%s\n",
			t.ToolCallerID, t.ToolName, t.Argument))
	}
}

func (l LLMMessage) TruncatedRawContent(placeholder string) string {
	return util.TruncateLines(l.RawContent, 3, 2, placeholder)
}

func TotalInputTokens(messages []LLMMessage) int64 {
	if len(messages) == 0 {
		return 0
	}

	total := int64(0)
	for _, msg := range messages {
		total += msg.Usage.TotalInputToken()
	}

	return total
}

func TotalOutputTokens(messages []LLMMessage) int64 {
	if len(messages) == 0 {
		return 0
	}

	total := int64(0)
	for _, msg := range messages {
		total += msg.Usage.TotalOutputToken()
	}

	return total
}

type ToolCall struct {
	ToolCallerID string
	ToolName     string
	Argument     string
}

type MessageRole string

const (
	LLMAssistant MessageRole = "assistant"
	LLMUser      MessageRole = "user"
	LLMSystem    MessageRole = "system"
	LLMTool      MessageRole = "tool"
)

type MessageFinishReason string

const (
	FinishStop       MessageFinishReason = "stop"
	FinishToolCalls  MessageFinishReason = "toolCalls"
	FinishLengthOver MessageFinishReason = "lengthOver"
)

type LLMUsage struct {
	InputToken       int64
	OutputToken      int64
	CacheCreateToken int64
	CacheReadToken   int64
}

func (l LLMUsage) TotalInputToken() int64 {
	return l.InputToken + l.CacheReadToken
}

func (l LLMUsage) TotalOutputToken() int64 {
	return l.OutputToken + l.CacheCreateToken
}
