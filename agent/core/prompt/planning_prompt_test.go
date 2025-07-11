package prompt_test

import (
	"testing"

	"github.com/clover0/issue-agent/core/prompt"
	"github.com/clover0/issue-agent/test/assert"
)

func TestPlanningPrompt_Build(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input prompt.PlanningPrompt
		want  prompt.Prompt
	}{
		"with all fields populated": {
			input: prompt.PlanningPrompt{
				Language:     "Japanese",
				BaseBranch:   "main",
				IssueTitle:   "Test Issue",
				IssueContent: "This is a test issue content",
				IssueNumber:  "123",
			},
			want: prompt.Prompt{
				SystemPrompt: `
You are a software development engineer with expertise in the latest technologies, programming, best practices.
You will be instructed by a user to accomplish the task.
Your goal is to plan what the developer needs to do to accomplish the task with instruction-format.
You must gather information and context from the repository to create a detailed plan for the developer.
To get information beyond the given task, you must read the files in the repository to get the information yourself.

The plan is passed as input to the agent doing the development.

<system-environment>
* You are in the root directory of the repository.
* Git Base branch is main
</system-environment>

<constraints>
* Communicate entirely in Japanese.
* You are in an environment where you cannot execute arbitrary commands, so you cannot run the shell. Only tool use can be used.
* Handling files with huge sizes is inefficient, so you can only open files that are less than 15,000 bytes.
* You and the developer work in the same environment.
</constraints>

<instruction-format>
* Specify the type of expert to act as (e.g., "You are an expert in developing applications using Go").
* Provide a step-by-step action plan of what needs to be done.
* Add a step to create a PR at the end of the plan.
</instruction-format>
`,
				StartUserPrompt: `
The task is bellow:

<task>
Title: Test Issue
This is a test issue content
</task>

<instructions>
* Think deeply about what is needed to complete the task.
* Thoroughly analyze the repository structure and source code to plan the development process.
* After planning, create instructions for the software development engineer who will complete the task.
* Finally, output only the instruction document written by English in the specified instruction-format for the software developer agent.
* Do not output anything other than the instruction document.
</instructions>
`,
			},
		},
		"with minimal fields": {
			input: prompt.PlanningPrompt{
				Language:     "English",
				BaseBranch:   "develop",
				IssueTitle:   "",
				IssueContent: "",
				IssueNumber:  "",
			},
			want: prompt.Prompt{
				SystemPrompt: `
You are a software development engineer with expertise in the latest technologies, programming, best practices.
You will be instructed by a user to accomplish the task.
Your goal is to plan what the developer needs to do to accomplish the task with instruction-format.
You must gather information and context from the repository to create a detailed plan for the developer.
To get information beyond the given task, you must read the files in the repository to get the information yourself.

The plan is passed as input to the agent doing the development.

<system-environment>
* You are in the root directory of the repository.
* Git Base branch is develop
</system-environment>

<constraints>
* Communicate entirely in English.
* You are in an environment where you cannot execute arbitrary commands, so you cannot run the shell. Only tool use can be used.
* Handling files with huge sizes is inefficient, so you can only open files that are less than 15,000 bytes.
* You and the developer work in the same environment.
</constraints>

<instruction-format>
* Specify the type of expert to act as (e.g., "You are an expert in developing applications using Go").
* Provide a step-by-step action plan of what needs to be done.
* Add a step to create a PR at the end of the plan.
</instruction-format>
`,
				StartUserPrompt: `
The task is bellow:

<task>
Title: 

</task>

<instructions>
* Think deeply about what is needed to complete the task.
* Thoroughly analyze the repository structure and source code to plan the development process.
* After planning, create instructions for the software development engineer who will complete the task.
* Finally, output only the instruction document written by English in the specified instruction-format for the software developer agent.
* Do not output anything other than the instruction document.
</instructions>
`,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.input.Build()
			assert.Nil(t, err)
			assert.Equal(t, got.SystemPrompt, tt.want.SystemPrompt)
			assert.Equal(t, got.StartUserPrompt, tt.want.StartUserPrompt)
		})
	}
}
