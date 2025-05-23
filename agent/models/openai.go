package models

// TODO: make no open-ai dependency
// The openai-go library is too large for the purposes of this project.

// TODO: move logic to communicate with LLM to the OpenAILLMForwarder struct

import (
	"context"
	"errors"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/clover0/issue-agent/core"
	"github.com/clover0/issue-agent/logger"
)

type OpenAI struct {
	client        *openai.Client
	sendLogger    logger.Logger
	receiveLogger logger.Logger
}

func NewOpenAI(lo logger.Logger, apiKey string) OpenAI {
	sendLogger := lo.AddPrefix("[OpenAIForwarder] ").SetColor(logger.Green)
	receiveLogger := lo.AddPrefix("[OpenAIReceive] ").SetColor(logger.Yellow)
	return OpenAI{
		sendLogger:    sendLogger,
		receiveLogger: receiveLogger,
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
	}
}

func (o OpenAI) createCompletionParams(input core.StartCompletionInput) (openai.ChatCompletionNewParams, []core.LLMMessage) {
	toolFuncs := make([]openai.ChatCompletionToolParam, len(input.Functions))
	for i, f := range input.Functions {
		toolFuncs[i] = openai.ChatCompletionToolParam{
			Function: openai.F(f.ToFunctionCalling()),
			Type:     openai.F(openai.ChatCompletionToolTypeFunction),
		}
	}

	historyInitial := []core.LLMMessage{
		{
			Role:       core.LLMSystem,
			RawContent: input.SystemPrompt,
		},
		{
			Role:       core.LLMUser,
			RawContent: input.StartUserPrompt,
		},
	}

	return openai.ChatCompletionNewParams{
		Model: openai.F(input.Model),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(input.SystemPrompt),
			openai.UserMessage(input.StartUserPrompt),
		}),
		Temperature: openai.F(0.0),
		Tools:       openai.F(toolFuncs),
	}, historyInitial
}

func (o OpenAI) StartCompletion(ctx context.Context, input core.StartCompletionInput) ([]core.LLMMessage, error) {
	var history []core.LLMMessage
	params, historyInitial := o.createCompletionParams(input)
	history = append(history, historyInitial...)

	o.sendLogger.Info(fmt.Sprintf("model: %s, sending message\n", input.Model))
	o.sendLogger.Debug("system prompt:\n%s\n", input.SystemPrompt)
	o.sendLogger.Debug("user prompt:\n%s\n", input.StartUserPrompt)
	chat, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}

	msg := chat.Choices[0]
	lastMsg := core.LLMMessage{
		Role:              core.LLMAssistant,
		RawContent:        msg.Message.Content,
		FinishReason:      convertToFinishReason(msg.FinishReason),
		ReturnedToolCalls: convertToToolCalls(msg.Message.ToolCalls),
		RawMessageStruct:  msg.Message,
		Usage: core.LLMUsage{
			InputToken:  int32(chat.Usage.PromptTokens),
			OutputToken: int32(chat.Usage.CompletionTokens),
			TotalToken:  int32(chat.Usage.TotalTokens),
		},
	}
	history = append(history, lastMsg)

	o.sendLogger.Debug(fmt.Sprintf("prompt token: %d, completion token: %d\n",
		chat.Usage.PromptTokens, chat.Usage.CompletionTokens,
	))

	o.receiveLogger.Info("returned messages:\n")
	lastMsg.ShowAssistantMessage(o.receiveLogger)

	return history, nil
}

func (o OpenAI) ContinueCompletion(
	ctx context.Context,
	input core.StartCompletionInput,
	llmContexts []core.ReturnToLLMContext,
	history []core.LLMMessage,
) ([]core.LLMMessage, error) {
	params, _ := o.createCompletionParams(input)

	// build from history
	params.Messages.Value = []openai.ChatCompletionMessageParamUnion{}
	for _, h := range history {
		switch h.Role {
		case core.LLMAssistant:
			if h.RawMessageStruct == nil {
				return nil, errors.New("rawMessageStruct should not be nil. But it is nil")
			}

			m, ok := h.RawMessageStruct.(openai.ChatCompletionMessage)
			if !ok {
				return nil, errors.New("RawMessageStruct can't convert ChatCompletionMessage")
			}

			params.Messages.Value = append(params.Messages.Value, m)
		case core.LLMUser:
			params.Messages.Value = append(params.Messages.Value, openai.UserMessage(h.RawContent))
		case core.LLMSystem:
			params.Messages.Value = append(params.Messages.Value, openai.SystemMessage(h.RawContent))
		case core.LLMTool:
			params.Messages.Value = append(params.Messages.Value,
				openai.ToolMessage(h.RespondToolCall.ToolCallerID, h.RawContent),
			)
		}
	}

	// new message
	var newMsg core.LLMMessage
	for _, v := range llmContexts {
		if v.ToolCallerID != "" {
			// tool message
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(v.ToolCallerID, v.Content))
			newMsg = core.LLMMessage{
				Role:       core.LLMTool,
				RawContent: v.Content,
				RespondToolCall: core.ToolCall{
					ToolCallerID: v.ToolCallerID,
					ToolName:     v.ToolName,
				},
			}
		} else {
			// user message
			params.Messages.Value = append(params.Messages.Value, openai.UserMessage(v.Content))
			newMsg = core.LLMMessage{
				Role:       core.LLMUser,
				RawContent: v.Content,
			}
		}
		history = append(history, newMsg)
	}

	o.debugShowSendingMsg(params)
	chat, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("continue completion error: %w", err)
	}

	msg := chat.Choices[0]
	lastMsg := core.LLMMessage{
		Role:              core.LLMAssistant,
		RawContent:        msg.Message.Content,
		FinishReason:      convertToFinishReason(msg.FinishReason),
		ReturnedToolCalls: convertToToolCalls(msg.Message.ToolCalls),
		RawMessageStruct:  msg.Message,
		Usage: core.LLMUsage{
			InputToken:  int32(chat.Usage.PromptTokens),
			OutputToken: int32(chat.Usage.CompletionTokens),
			TotalToken:  int32(chat.Usage.TotalTokens),
		},
	}
	history = append(history, lastMsg)

	o.receiveLogger.Info("returned messages:\n")
	lastMsg.ShowAssistantMessage(o.receiveLogger)

	return history, nil
}

func convertToFinishReason(finishReason openai.ChatCompletionChoicesFinishReason) core.MessageFinishReason {
	switch finishReason {
	case openai.ChatCompletionChoicesFinishReasonLength:
		return core.FinishLengthOver
	case openai.ChatCompletionChoicesFinishReasonStop:
		return core.FinishStop
	case openai.ChatCompletionChoicesFinishReasonToolCalls:
		return core.FinishToolCalls
	default:
		panic(fmt.Sprintf("convertToFinishReason: unknown finish reason: %s", finishReason))
	}
}

func convertToToolCalls(toolCalls []openai.ChatCompletionMessageToolCall) []core.ToolCall {
	var res []core.ToolCall
	for _, v := range toolCalls {
		res = append(res, core.ToolCall{
			ToolCallerID: v.ID,
			ToolName:     v.Function.Name,
			Argument:     v.Function.Arguments,
		})
	}
	return res
}

func (o OpenAI) CompletionNextStep(_ context.Context, history []core.LLMMessage) core.Step {
	// last message
	lastMsg := history[len(history)-1]

	switch lastMsg.FinishReason {
	case core.FinishStop:
		return core.NewWaitingInstructionStep(lastMsg.RawContent)
	case core.FinishToolCalls:
		var input []core.FunctionsInput
		for _, v := range lastMsg.ReturnedToolCalls {
			input = append(input, core.FunctionsInput{
				FuncName:     v.ToolName,
				FunctionArgs: v.Argument,
				ToolCallerID: v.ToolCallerID,
			})
		}
		return core.NewExecStep(input)
	case core.FinishLengthOver:
		return core.NewUnrecoverableStep(fmt.Errorf("chat completion length error"))
	}

	return core.NewUnknownStep()
}

func (o OpenAI) debugShowSendingMsg(param openai.ChatCompletionNewParams) {
	if len(param.Messages.Value) > 0 {
		o.sendLogger.Info(fmt.Sprintf("model: %s, sending messages:\n", param.Model.String()))
		// TODO: show all messages. But now, show only the last message
		o.sendLogger.Debug(fmt.Sprintf("%s\n", param.Messages.Value[len(param.Messages.Value)-1]))
	}
}
