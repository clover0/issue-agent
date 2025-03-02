package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	"github.com/clover0/issue-agent/agent"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/util/pointer"
)

// TODO: refactor using ptr package

type BedrockLLMForwarder struct {
	Bedrock BedrockClient
}

func NewBedrockLLMForwarder(l logger.Logger) (agent.LLMForwarder, error) {
	bed, err := NewBedrock(l)
	if err != nil {
		return nil, err
	}
	return BedrockLLMForwarder{
		Bedrock: bed,
	}, nil
}

func (a BedrockLLMForwarder) StartForward(input agent.StartCompletionInput) ([]agent.LLMMessage, error) {
	var history []agent.LLMMessage
	initMsg, toolSpecs, initialHistory := a.buildStartParams(input)

	history = append(history, initialHistory...)

	a.Bedrock.logger.Info(logger.Green(fmt.Sprintf("model: %s, sending message\n", input.Model)))
	a.Bedrock.logger.Debug("system prompt:\n%s\n", input.SystemPrompt)
	a.Bedrock.logger.Debug("user prompt:\n%s\n", input.StartUserPrompt)
	resp, err := a.Bedrock.Messages.Create(context.TODO(),
		input.Model,
		input.SystemPrompt,
		initMsg,
		toolSpecs,
	)
	if err != nil {
		return nil, err
	}

	assistantHist, err := a.buildAssistantHistory(*resp)
	if err != nil {
		return nil, err
	}

	history = append(history, assistantHist)

	a.Bedrock.logger.Info(logger.Yellow("returned messages:\n"))
	assistantHist.ShowAssistantMessage(a.Bedrock.logger)

	return history, nil
}

func (a BedrockLLMForwarder) ForwardLLM(
	_ context.Context,
	input agent.StartCompletionInput,
	llmContexts []agent.ReturnToLLMContext,
	history []agent.LLMMessage,
) ([]agent.LLMMessage, error) {
	_, toolSpecs, _ := a.buildStartParams(input)

	// reset message
	var messages []types.Message

	// build message from history
	for _, h := range history {
		var msg types.Message

		switch h.Role {
		case agent.LLMAssistant:
			msg.Role = types.ConversationRoleAssistant
			if len(h.ReturnedToolCalls) > 0 {
				for _, v := range h.ReturnedToolCalls {
					var inputMap map[string]any
					if err := json.Unmarshal([]byte(v.Argument), &inputMap); err != nil {
						return nil, fmt.Errorf("failed to unmarshal tool argument: %w", err)
					}
					msg.Content = append(msg.Content, &types.ContentBlockMemberToolUse{
						Value: types.ToolUseBlock{
							ToolUseId: &v.ToolCallerID,
							Name:      &v.ToolName,
							Input:     document.NewLazyDocument(inputMap),
						},
					})
				}
			} else {
				msg.Content = append(msg.Content, &types.ContentBlockMemberText{
					Value: h.RawContent,
				})
			}

		case agent.LLMUser:
			msg.Role = types.ConversationRoleUser
			msg.Content = append(msg.Content, &types.ContentBlockMemberText{
				Value: h.RawContent,
			})

		case agent.LLMTool:
			msg.Role = types.ConversationRoleUser
			toolResult := types.ToolResultBlock{
				ToolUseId: &h.RespondToolCall.ToolCallerID,
			}
			toolResult.Content = append(toolResult.Content, &types.ToolResultContentBlockMemberText{Value: h.RawContent})
			msg.Content = append(msg.Content, &types.ContentBlockMemberToolResult{Value: toolResult})

		default:
			return nil, fmt.Errorf("unknown role: %s", h.Role)
		}

		messages = append(messages, msg)
	}

	// new message
	var newMsg agent.LLMMessage
	content := make([]types.ContentBlock, len(llmContexts))
	for i, v := range llmContexts {
		// only one content ?
		if v.ToolCallerID != "" {
			content[i] = &types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					ToolUseId: &v.ToolCallerID,
					Content: []types.ToolResultContentBlock{
						&types.ToolResultContentBlockMemberText{
							Value: v.Content,
						},
					},
				},
			}
			newMsg = agent.LLMMessage{
				Role:       agent.LLMTool,
				RawContent: v.Content,
				RespondToolCall: agent.ToolCall{
					ToolCallerID: v.ToolCallerID,
					ToolName:     v.ToolName,
				},
			}
		} else {
			content[i] = &types.ContentBlockMemberText{
				Value: v.Content,
			}
			newMsg = agent.LLMMessage{
				Role:       agent.LLMUser,
				RawContent: v.Content,
			}
		}

		history = append(history, newMsg)
	}

	messages = append(messages, types.Message{
		Role:    types.ConversationRoleUser,
		Content: content,
	})

	a.Bedrock.logger.Info(logger.Green(fmt.Sprintf("model: %s, sending message\n", input.Model)))
	a.Bedrock.logger.Debug("%s\n", newMsg.RawContent)

	resp, err := a.Bedrock.Messages.Create(
		context.TODO(),
		input.Model,
		input.SystemPrompt,
		messages,
		toolSpecs)
	if err != nil {
		return nil, err
	}

	assistantHist, err := a.buildAssistantHistory(*resp)
	if err != nil {
		return nil, err
	}
	history = append(history, assistantHist)

	a.Bedrock.logger.Info(logger.Yellow("LLM returned messages:\n"))
	assistantHist.ShowAssistantMessage(a.Bedrock.logger)

	return history, nil
}

// TODO: refactor with openai forwarder
func (a BedrockLLMForwarder) ForwardStep(_ context.Context, history []agent.LLMMessage) agent.Step {
	lastMsg := history[len(history)-1]

	switch lastMsg.FinishReason {
	case agent.FinishStop:
		return agent.NewWaitingInstructionStep(lastMsg.RawContent)
	case agent.FinishToolCalls:
		var input []agent.FunctionsInput
		for _, v := range lastMsg.ReturnedToolCalls {
			input = append(input, agent.FunctionsInput{
				FuncName:     v.ToolName,
				FunctionArgs: v.Argument,
				ToolCallerID: v.ToolCallerID,
			})
		}
		return agent.NewExecStep(input)
	case agent.FinishLengthOver:
		return agent.NewUnrecoverableStep(fmt.Errorf("chat completion length error"))
	}

	return agent.NewUnknownStep()
}

func (a BedrockLLMForwarder) buildAssistantHistory(bedrockResp bedrockruntime.ConverseOutput) (agent.LLMMessage, error) {
	respMessage, ok := bedrockResp.Output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return agent.LLMMessage{}, fmt.Errorf("failed to convert output to message")
	}
	var toolCalls []agent.ToolCall
	var text string
	for _, cont := range respMessage.Value.Content {
		switch c := cont.(type) {
		case *types.ContentBlockMemberText:
			text = c.Value
		case *types.ContentBlockMemberToolUse:
			doc, err := c.Value.Input.MarshalSmithyDocument()
			if err != nil {
				return agent.LLMMessage{}, fmt.Errorf("failed to unmarshal tool argument: %w", err)
			}
			toolCalls = append(toolCalls, agent.ToolCall{
				ToolCallerID: *c.Value.ToolUseId,
				ToolName:     *c.Value.Name,
				Argument:     string(doc),
			})
		default:
			return agent.LLMMessage{}, fmt.Errorf("unknown content type: %T", c)
		}
	}

	return agent.LLMMessage{
		Role:              agent.LLMAssistant,
		FinishReason:      convertBedrockStopReasonToReason(bedrockResp.StopReason),
		RawContent:        text,
		ReturnedToolCalls: toolCalls,
		Usage: agent.LLMUsage{
			InputToken:  *bedrockResp.Usage.InputTokens,
			OutputToken: *bedrockResp.Usage.OutputTokens,
			TotalToken:  *bedrockResp.Usage.TotalTokens,
		},
	}, nil
}

// TODO: refactor rename
func (a BedrockLLMForwarder) buildStartParams(input agent.StartCompletionInput) ([]types.Message, []*types.ToolMemberToolSpec, []agent.LLMMessage) {
	var messages []types.Message
	tools := make([]*types.ToolMemberToolSpec, len(input.Functions))

	for i, f := range input.Functions {
		tools[i] = &types.ToolMemberToolSpec{
			Value: types.ToolSpecification{
				Name:        pointer.String(f.Name.String()),
				Description: &f.Description,
				InputSchema: &types.ToolInputSchemaMemberJson{
					Value: document.NewLazyDocument(f.Parameters),
				},
			},
		}
	}

	messages = append(messages, types.Message{
		Role: types.ConversationRoleUser,
		Content: []types.ContentBlock{
			&types.ContentBlockMemberText{
				Value: input.StartUserPrompt,
			},
		},
	})

	return messages, tools, []agent.LLMMessage{
		{
			Role:       agent.LLMUser,
			RawContent: input.StartUserPrompt,
		},
	}
}

func convertBedrockStopReasonToReason(reason types.StopReason) agent.MessageFinishReason {
	switch reason {
	case types.StopReasonEndTurn:
		return agent.FinishStop
	case types.StopReasonToolUse:
		return agent.FinishToolCalls
	case types.StopReasonMaxTokens:
		return agent.FinishLengthOver
	default:
		panic(fmt.Sprintf("unknown finish reason: %s", reason))
	}
}
