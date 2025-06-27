package functions

import (
	"strings"
)

const FuncInvokeAgent = "invoke_agent"

type InvokeAgentType func(input InvokeAgentInput) (InvokeAgentOutput, error)

type AgentInvokerIF interface {
	Invoke(input InvokeAgentInput) (InvokeAgentOutput, error)
}

func InitInvokeAgentFunction(agentInvoker AgentInvokerIF) Function {
	f := Function{
		Name: FuncInvokeAgent,
		Description: strings.ReplaceAll(`Run an AI Agent powered by LLM with your system prompt and first user prompt.
AI Agents require relevant context to function properly. While the Git environment is shared, other contextual information must be provided externally
through mechanisms such as system prompts and first user prompts.
This includes information such as the current branch, Pull Request number tnd issue number.
When completing work, it is essential to output what was accomplished so that other AI agents can understand what was done.`,
			"\n", " "),
		Func: InvokeAgentCaller(agentInvoker),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "The name of the agent.",
					"minLength":   3,
					"maxLength":   20,
				},
				"system_prompt": map[string]any{
					"type": "string",
					"description": strings.ReplaceAll(
						`System prompt is an instruction given to AI systems that define their behavior parameters, including role,
response style, and functional limitations, set invisibly before conversations begin.`,
						"\n", " "),
				},
				"first_user_prompt": map[string]any{
					"type":        "string",
					"description": "The first user prompt is the initial question or command given to the AI agent.",
				},
			},
			"required":             []string{"name", "system_prompt", "first_user_prompt"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncInvokeAgent] = f

	return f
}

type InvokeAgentInput struct {
	Name            string `json:"name"`
	SystemPrompt    string `json:"system_prompt"`
	FirstUserPrompt string `json:"first_user_prompt"`
}

type InvokeAgentOutput struct {
	Content string
}

func (r InvokeAgentOutput) ToLLMString() string {
	return r.Content
}

func InvokeAgentCaller(
	agentInvoker AgentInvokerIF,
) InvokeAgentType {
	return func(input InvokeAgentInput) (InvokeAgentOutput, error) {
		return agentInvoker.Invoke(input)
	}
}
