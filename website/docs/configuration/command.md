# Command

```shell
$ issue-agent help

Usage
  issue-agent <command> [flags]
Command and Flags
  help: Show usage of commands and flags
  version: Show version of issue-agent CLI
  create-pr:
    Usage:
      create-pr GITHUB_OWNER/REPOSITORY/issues/NUMBER [flags]
    Flags:
    --aws_profile
      AWS profile to use a specific profile from credentials.
    --aws_region
      AWS region to use for credentials and Bedrock.
      Default(If use aws_profile): aws profile's default session region.
    --base_branch
      Base Branch for pull request
    --config
      Path to the configuration file.
      Default: agent/config/default_config.yml in this project.
    --language
      Language spoken by agent.
      Default: English.
    --log_level
      Log level. If you want to see LLM completions, set it to 'debug'.
      Default: info.
    --model
      LLM name. For the model name, check the documentation of each LLM provider.
    --review_agents
      The number of agents to review.
      If a value greater than 0 is specified, then that number of reviews will be performed by agents with different roles.
      Default: 0

```
