<p align="center">
  <h1 align="center">Issue Agent</h1>
  <p align="center">An AI Agent that quickly solves simple issues</p>
</p>

---


The Issue Agent is a tool that generates work products in a Git repository based on a given assignment.

Powered by Large Language Models (LLMs), 
this agent is designed to help developers efficiently tackle simple issues. 
When a developer creates an issue in a repository and assigns it to this agent, 
it autonomously works to solve the issue and submits the results as a Pull Request on GitHub.


## Install

- [Your Machine](https://clover0.github.io/issue-agent/getting-started/installation/)
- [GitHub Action](https://github.com/clover0/setup-issue-agent)


## Documentation
Refer to the [documentation](https://clover0.github.io/issue-agent) for more details.


## Supported Models
The following models are supported.

- OpenAI Models
  - gpt-4o
  - gpt-4o-mini
- Anthropic Models
  - claude-3-5-sonnet-latest (⭐️Strongly Recommended!)
- AWS Bedrock Models
  - claude-3-5-sonnet v2 (ModelID = anthropic.claude-3-5-sonnet-20241022-v2:0)
  - claude-3-5-sonnet v2 (ModelID = us.anthropic.claude-3-5-sonnet-20241022-v2:0, Cross-region inference)
  - claude-3-5-sonnet v1 (ModelID = anthropic.claude-3-5-sonnet-20240620-v1:0)
  - claude-3-5-sonnet v1 (ModelID = us.anthropic.claude-3-5-sonnet-20240620-v1:0, Cross-region inference)
