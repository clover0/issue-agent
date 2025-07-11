# Command

```
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

  react:
    Usage:
      react RESOURCE_FORMAT [flags]
    RESOURCE_FORMAT:
        issue_comment(pull request comment): OWNER/REPO/issues/comments/COMMENT_ID
        pull_request_review_comment: OWNER/REPO/pulls/comments/COMMENT_ID
    Example:
       react owner/example/issues/comments/123456 [flags]
    Flags:
    --aws_profile
      AWS profile to use a specific profile from credentials.
    --aws_region
      AWS region to use for credentials and Bedrock.
      Default(If use aws_profile): aws profile's default session region.
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

```


## `create-pr` command

Using functions:

- get_pull_request
- list_files
- modify_file
- open_file
- put_file
- submit_files
- search_files
- remove_file
- switch_branch
- submit_revision
- get_issue
- create_pull_request_comment
- get_repository_content
- invoke_agent
- request_reviewers


## `react` command

Using functions:

- get_pull_request
- list_files
- modify_file
- open_file
- put_file
- search_files
- remove_file
- submit_revision
- get_issue
- create_pull_request_comment
- create_pull_request_review_comment
- get_repository_content


Issue Agent does not save prompt history.
Therefore, When user uses the `react` command, the agent will not remember the previous conversation.
