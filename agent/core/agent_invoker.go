package core

import (
	"fmt"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/core/prompt"
	"github.com/clover0/issue-agent/logger"
)

type AgentInvoker struct {
	params    Parameter
	logg      logger.Logger
	forwarder LLMForwarder
	tools     []functions.Function
}

func NewAgentInvoker(
	params Parameter,
	logg logger.Logger,
	forwarder LLMForwarder,
	tools []functions.Function,
) functions.AgentInvokerIF {
	return &AgentInvoker{
		params:    params,
		logg:      logg,
		forwarder: forwarder,
		tools:     tools,
	}
}

func (a AgentInvoker) Invoke(input functions.InvokeAgentInput) (functions.InvokeAgentOutput, error) {
	if input.Name == "" {
		return functions.InvokeAgentOutput{}, fmt.Errorf("name is required")
	}
	if input.SystemPrompt == "" {
		return functions.InvokeAgentOutput{}, fmt.Errorf("system_prompt is required")
	}
	if input.FirstUserPrompt == "" {
		return functions.InvokeAgentOutput{}, fmt.Errorf("first_user_prompt is required")
	}

	l := a.logg.AddPrefix(fmt.Sprintf("[agent][%s]", input.Name))
	agent := NewAgent(
		a.params,
		input.Name,
		l,
		prompt.Prompt{
			SystemPrompt:    input.SystemPrompt,
			StartUserPrompt: input.FirstUserPrompt,
		},
		a.forwarder,
		a.tools,
	)

	lastOutput, err := agent.Work()
	if err != nil {
		return functions.InvokeAgentOutput{}, fmt.Errorf("failed to run agent: %w", err)
	}

	return functions.InvokeAgentOutput{Content: lastOutput}, nil
}
