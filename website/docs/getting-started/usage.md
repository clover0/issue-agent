# Usage

## Run Standard
Choose the GitHub issue you want to resolve.

### Example

- GitHub Repository: `clover0/example-repository`
- Issue Number: 123 to solve
- Base Branch: `main` to create a pull request
- LLM: Anthropic Claude 3.5 Sonnet

```shell
$ issue-agent create-pr clover0/example-repository/issues/123 \
  --base_branch main \
  --model claude-3-5-sonnet-latest \
```

## Run with Environment Variables
You can use environment variables to run the `issue-agent` in a single line.

The [`gh` CLI](https://github.com/cli/cli#installation) is particularly useful for managing tokens.

```shell
$ GITHUB_TOKEN=$(gh auth token) \
  ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY} \
  issue-agent create-pr clover0/example-repository/issues/123 \
    --base_branch your_branch  \
    --model claude-3-5-sonnet-latest \
```


## Run AWS Bedrock with SSO session
You can also execute the `issue-agent` command using AWS Bedrock with an SSO session.

```sh
$ GITHUB_TOKEN=$(gh auth token) \
issue-agent create-pr clover0/example-repository/issues/123 \
  --base_branch your_branch  \
  --model anthropic.claude-3-5-sonnet-20241022-v2:0 \
  --aws_profile your_profile \
  --aws_region us-east-1
```


## Branch Naming

The working branch is created automatically with a prefix of `agent-`.
