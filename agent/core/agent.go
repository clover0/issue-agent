package core

import (
	"context"
	"fmt"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/core/prompt"
	"github.com/clover0/issue-agent/core/store"
	"github.com/clover0/issue-agent/logger"
)

type AgentLike interface {
	Work() (lastOutput string, err error)
	History() []LLMMessage
	LastHistory() LLMMessage
	ChangedFiles() []store.File
}

type Agent struct {
	name         string
	parameter    Parameter
	currentStep  Step
	logg         logger.Logger
	llmForwarder LLMForwarder
	prompt       prompt.Prompt
	history      []LLMMessage
	store        *store.Store
	tools        []functions.Function
}

func NewAgent(
	parameter Parameter,
	name string,
	logg logger.Logger,
	prompt prompt.Prompt,
	forwarder LLMForwarder,
	store *store.Store,
	tools []functions.Function,
) AgentLike {
	return &Agent{
		name:         name,
		parameter:    parameter,
		currentStep:  Step{},
		logg:         logg,
		prompt:       prompt,
		llmForwarder: forwarder,
		store:        store,
		tools:        tools,
	}
}

func (a *Agent) Work() (lastOutput string, err error) {
	ctx := context.Background()
	a.logg.Info("[%s]agent starts work\n", a.name)

	completionInput := StartCompletionInput{
		Model:           a.parameter.Model,
		SystemPrompt:    a.prompt.SystemPrompt,
		StartUserPrompt: a.prompt.StartUserPrompt,
		Functions:       a.tools,
	}

	logGreen, logBlue, logRed := a.logg.SetColor(logger.Green), a.logg.SetColor(logger.Blue), a.logg.SetColor(logger.Red)
	logGreen.Info("[STEP:1]start communication with LLM\n")
	history, err := a.llmForwarder.StartForward(completionInput)
	if err != nil {
		return lastOutput, fmt.Errorf("start llm forward error: %w", err)
	}
	a.updateHistory(history)

	a.currentStep = a.llmForwarder.ForwardStep(ctx, history)

	var steps = 1
	loop := true
	for loop {
		steps++
		if steps > a.parameter.MaxSteps {
			a.logg.Info(fmt.Sprintf("reached to the max steps %d\n", a.parameter.MaxSteps))
			break
		}
		stepLabel := fmt.Sprintf("[STEP:%d]", steps)

		switch a.currentStep.Do {
		case Exec:
			logBlue.Info(stepLabel + "execute functions:\n")
			var input []ReturnToLLMInput
			for _, fnCtx := range a.currentStep.FunctionContexts {
				var returningStr string
				returningStr, err = functions.ExecFunction(
					a.logg,
					a.store,
					fnCtx.Function.Name,
					fnCtx.FunctionArgs.String(),
				)

				if err != nil {
					logRed.Error("function error"+": %s\n", err)
					returningStr = fmt.Sprintf("Error caused. error message: %s\nChange the arguments before using it again. "+
						"If you still get an error, change the tool you are using", err.Error())
				}

				input = append(input, ReturnToLLMInput{
					ToolCallerID: fnCtx.ToolCallerID,
					Content:      returningStr,
				})
			}
			a.currentStep = NewReturnToLLMStep(input)

		case ReturnToLLM:
			logGreen.Info(stepLabel + "forwarding message to LLM and waiting for response\n")
			history, err = a.llmForwarder.ForwardLLM(ctx, completionInput, a.currentStep.ReturnToLLMContexts, history)
			if err != nil {
				a.logg.Error("unrecoverable error: %s\n", err)
				return lastOutput, err
			}
			a.updateHistory(history)
			a.currentStep = a.llmForwarder.ForwardStep(ctx, history)

		case WaitingInstruction:
			a.logg.Info(stepLabel + "finish instructions\n")
			lastOutput = a.currentStep.LastOutput
			loop = false

		case Unrecoverable, Unknown:
			a.logg.Error("unrecoverable error: %s\n", a.currentStep.UnrecoverableErr)
			return lastOutput, fmt.Errorf("unrecoverable error: %s", a.currentStep.UnrecoverableErr)

		default:
			a.logg.Error("%s does not exist in step types\n", a.currentStep.Do)
			return lastOutput, fmt.Errorf("%s does not exist in step type", a.currentStep.Do)
		}
	}
	return lastOutput, nil
}

func (a *Agent) updateHistory(history []LLMMessage) {
	a.history = history
}

func (a *Agent) History() []LLMMessage {
	return a.history
}

func (a *Agent) LastHistory() LLMMessage {
	return a.history[len(a.history)-1]
}

func (a *Agent) ChangedFiles() []store.File {
	return a.store.ChangedFiles()
}
