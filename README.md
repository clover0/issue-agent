<p align="center">
  <h1 align="center">Issue Agent</h1>
  <p align="center">An AI Agent that quickly solves simple GitHub issues</p>
</p>

---

Powered by Large Language Models (LLMs).

When a developer creates an issue in a repository and passes to this agent, 
it autonomously works to solve the issue and submits its results as a Pull Request on GitHub.


## Quick Installation & Usage
### Installation
- [Your Machine](https://clover0.github.io/issue-agent/getting-started/installation/)
- [GitHub Action](https://github.com/clover0/setup-issue-agent)


### Usage
```shell
issue-agent help
```

```shell
$ issue-agent create-pr clover0/example-repository/issues/123 \
  --base_branch main \
  --model claude-3-7-sonnet-20250219
```


## Documentation
Refer to the [documentation](https://clover0.github.io/issue-agent) for more details.


## Key Features
- **Fully Autonomous**
  - Handles simple coding and documentation tasks without human intervention.

- **Minimal Configuration**
  - Easy to set up: runs locally or via GitHub Actions in minutes.

- **Security**
  - Limited scope: cannot execute arbitrary or unsafe code (only predefined and controlled functions).
  - Never leak credentials or secrets.


## Suitable Use Cases
Issue Agent makes life easier, especially for routine or repetitive tasks:

- **Routine Development Tasks**
  _(e.g., basic migration, formatting, simple refactoring tasks)_
- **Documentation Maintenance**
  _(generate or update consistently formatted docs)_
- **Code Cleanup**
  _(e.g., batch updates, typo corrections, simple bug fixes)_


## Supported AI Models
We recommend Anthropic's Claude models for optimal performance:

| Provider          | Supported Models                                     |
|-------------------|------------------------------------------------------|
| **OpenAI**        | gpt-4o, gpt-4o-mini                                  |
| **Anthropic**⭐️   | claude-3-5-sonnet-latest, claude-3-7-sonnet-20250219 |
| **AWS Bedrock**⭐️ | See AWS Bedrock section                              |


### AWS Bedrock
The following models are supported.

- claude-3-5-sonnet v2 (ModelID = anthropic.claude-3-5-sonnet-20241022-v2:0)
- claude-3-5-sonnet v2 (ModelID = us.anthropic.claude-3-5-sonnet-20241022-v2:0, Cross-region inference)
- claude-3-5-sonnet v1 (ModelID = anthropic.claude-3-5-sonnet-20240620-v1:0)
- claude-3-5-sonnet v1 (ModelID = us.anthropic.claude-3-5-sonnet-20240620-v1:0, Cross-region inference)
- claude-3-7-sonnet (ModelID = anthropic.claude-3-7-sonnet-20250219-v1:0)
- claude-3-7-sonnet (ModelID = us.anthropic.claude-3-7-sonnet-20250219-v1:0, Cross-region inference)
