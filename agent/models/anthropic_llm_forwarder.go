package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/clover0/issue-agent/core"
	"github.com/clover0/issue-agent/logger"
)

type AnthropicLLMForwarder struct {
	anthropic     AnthropicClient
	forwardLogger logger.Logger
	receiveLogger logger.Logger
}

func NewAnthropicLLMForwarder(l logger.Logger) (core.LLMForwarder, error) {
	token, ok := os.LookupEnv("ANTHROPIC_API_KEY")
	if !ok {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	forwardLogger := l.AddPrefix("[AnthropicForwarder] ").SetColor(logger.Green)
	receiveLogger := l.AddPrefix("[AnthropicReceive] ").SetColor(logger.Yellow)

	return AnthropicLLMForwarder{
		anthropic:     NewAnthropic(l, token),
		forwardLogger: forwardLogger,
		receiveLogger: receiveLogger,
	}, nil
}

func (a AnthropicLLMForwarder) StartForward(input core.StartCompletionInput) ([]core.LLMMessage, error) {
	var history []core.LLMMessage
	params, initialHistory := a.createParams(input)
	history = append(history, initialHistory...)

	a.forwardLogger.Info(fmt.Sprintf("model: %s, sending message\n", input.Model))
	a.forwardLogger.Debug("system prompt:\n%s\n", input.SystemPrompt)
	a.forwardLogger.Debug("user prompt:\n%s\n", input.StartUserPrompt)
	resp, err := a.anthropic.Messages.Create(context.TODO(), params)
	if err != nil {
		return nil, err
	}

	var toolCalls []core.ToolCall
	var text string
	for _, cont := range resp.Content {
		// discard text
		if cont.Type == "text" {
			text = cont.Text
			continue
		}
		if cont.Type == "tool_use" {
			j, err := json.Marshal(cont.Input)
			if err != nil {
				return nil, err
			}
			toolCalls = append(toolCalls, core.ToolCall{
				ToolCallerID: cont.ID,
				ToolName:     cont.Name,
				Argument:     string(j),
			})
		}
	}
	lastMsg := core.LLMMessage{
		Role:              core.LLMAssistant,
		FinishReason:      convertAnthropicStopReasonToReason(resp.StopReason),
		RawContent:        text,
		ReturnedToolCalls: toolCalls,
		Usage: core.LLMUsage{
			InputToken:       resp.Usage.InputTokens,
			OutputToken:      resp.Usage.OutputTokens,
			CacheReadToken:   resp.Usage.CacheReadInputTokens,
			CacheCreateToken: resp.Usage.CacheCreationInputTokens,
		},
	}
	history = append(history, lastMsg)

	a.receiveLogger.Info("returned messages:\n")
	lastMsg.ShowAssistantMessage(a.anthropic.logger)

	return history, nil
}

func (a AnthropicLLMForwarder) ForwardLLM(
	_ context.Context,
	input core.StartCompletionInput,
	llmContexts []core.ReturnToLLMContext,
	history []core.LLMMessage,
) ([]core.LLMMessage, error) {
	params, _ := a.createParams(input)

	// reset message
	params["messages"] = make([]J, 0)

	// build message from history
	for _, h := range history {
		switch h.Role {
		case core.LLMAssistant:
			if len(h.ReturnedToolCalls) > 0 {
				content := make([]J, 0)
				for _, v := range h.ReturnedToolCalls {
					var input map[string]any
					if err := json.Unmarshal([]byte(v.Argument), &input); err != nil {
						return nil, fmt.Errorf("failed to unmarshal tool argument: %w", err)
					}
					content = append(content, J{
						"type": "tool_use",
						"id":   v.ToolCallerID,
						"name": v.ToolName,
						// json marshal?
						"input": input,
					})
				}
				params["messages"] = append(params["messages"].([]J), J{
					"role":    "assistant",
					"content": content,
				})
			} else {
				params["messages"] = append(params["messages"].([]J), J{
					"role":    "assistant",
					"content": h.RawContent,
				})
			}
		case core.LLMUser:
			params["messages"] = append(params["messages"].([]J), J{
				"role":    "user",
				"content": h.RawContent,
			})

		// multiple contents in 1 message
		case core.LLMTool:
			// Multiple LLM messages should be consolidated into a single message.
			params["messages"] = append(params["messages"].([]J), J{
				"role": "user",
				"content": []J{
					{
						"type":        "tool_result",
						"tool_use_id": h.RespondToolCall.ToolCallerID,
						"content":     h.RawContent,
					},
				},
			})
		default:
			return nil, fmt.Errorf("unknown role: %s", h.Role)
		}
	}

	// new message
	var newMsg core.LLMMessage
	content := make([]J, len(llmContexts))
	for i, v := range llmContexts {
		if v.ToolCallerID != "" {
			contentContent := []J{
				{
					"type": "text",
					"text": v.Content,
				},
			}

			// cache_control cannot be set for empty text blocks
			if v.Content != "" {
				contentContent[0]["cache_control"] = J{"type": "ephemeral"}
			}

			content[i] = J{
				"type":        "tool_result",
				"tool_use_id": v.ToolCallerID,
				"content":     contentContent,
			}

			newMsg = core.LLMMessage{
				Role:       core.LLMTool,
				RawContent: v.Content,
				RespondToolCall: core.ToolCall{
					ToolCallerID: v.ToolCallerID,
					ToolName:     v.ToolName,
				},
			}
		} else {
			params["messages"] = append(params["messages"].([]J), J{
				"role": "user",
				"content": []J{
					{
						"type":          "text",
						"text":          v.Content,
						"cache_control": J{"type": "ephemeral"},
					},
				},
			})
			newMsg = core.LLMMessage{
				Role:       core.LLMUser,
				RawContent: v.Content,
			}
		}
		history = append(history, newMsg)
	}
	params["messages"] = append(params["messages"].([]J), J{
		"role":    "user",
		"content": content,
	})

	a.forwardLogger.Info(fmt.Sprintf("model: %s, sending message\n", input.Model))
	a.forwardLogger.Debug("%s\n", newMsg.TruncatedRawContent("... truncated in debug output"))

	resp, err := a.anthropic.Messages.Create(context.TODO(), params)
	if err != nil {
		return nil, err
	}

	// TODO: refactor with StartForward
	var toolCalls []core.ToolCall
	var text string
	for _, cont := range resp.Content {
		// assumption of only 1 text per content
		if cont.Type == "text" {
			text = cont.Text
			continue
		}
		if cont.Type == "tool_use" {
			j, err := json.Marshal(cont.Input)
			if err != nil {
				return nil, err
			}
			toolCalls = append(toolCalls, core.ToolCall{
				ToolCallerID: cont.ID,
				ToolName:     cont.Name,
				Argument:     string(j),
			})
		}
	}

	lastMsg := core.LLMMessage{
		Role:              core.LLMAssistant,
		FinishReason:      convertAnthropicStopReasonToReason(resp.StopReason),
		RawContent:        text,
		ReturnedToolCalls: toolCalls,
		Usage: core.LLMUsage{
			InputToken:       resp.Usage.InputTokens,
			OutputToken:      resp.Usage.OutputTokens,
			CacheReadToken:   resp.Usage.CacheReadInputTokens,
			CacheCreateToken: resp.Usage.CacheCreationInputTokens,
		},
	}
	history = append(history, lastMsg)

	a.receiveLogger.Info("returned messages:\n")
	lastMsg.ShowAssistantMessage(a.anthropic.logger)

	return history, nil
}

// TODO: refactor with openai forwarder
func (a AnthropicLLMForwarder) ForwardStep(_ context.Context, history []core.LLMMessage) core.Step {
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

func (a AnthropicLLMForwarder) createParams(input core.StartCompletionInput) (J, []core.LLMMessage) {
	tools := make([]J, len(input.Functions))

	for i, f := range input.Functions {
		tools[i] = J{
			"name":         f.Name,
			"description":  f.Description,
			"input_schema": f.Parameters,
		}
	}

	body := J{
		"model": input.Model,
		"system": []J{
			{
				"type":          "text",
				"text":          input.SystemPrompt,
				"cache_control": J{"type": "ephemeral"},
			},
		},
		"messages": []J{
			{"role": "user", "content": []J{
				{
					"type":          "text",
					"text":          input.StartUserPrompt,
					"cache_control": J{"type": "ephemeral"},
				},
			}},
		},
		"temperature": 0.0,
		"tool_choice": J{
			"type":                      "auto",
			"disable_parallel_tool_use": true,
		},
		"tools":      tools,
		"max_tokens": ClaudeMaxOutputTokens(input.Model),
	}

	return body, []core.LLMMessage{
		{
			Role:       core.LLMUser,
			RawContent: input.StartUserPrompt,
		},
	}
}

// TODO: refactor to shared multi models
func convertAnthropicStopReasonToReason(reason string) core.MessageFinishReason {
	switch reason {
	case "end_turn":
		return core.FinishStop
	case "max_tokens":
		return core.FinishLengthOver
	case "stop_sequence":
		return core.FinishStop
	case "too_use":
		return core.FinishToolCalls
	default:
		return core.FinishToolCalls
	}
}
