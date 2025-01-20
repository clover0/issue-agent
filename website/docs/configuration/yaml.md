# Configuration YAML

At your repository root, create a `issue_agent.yml` file with the following content.

```yaml
# Communication language
# English, Japanese...
# Default is English
language: "English"

# Default is /tmp/repositories
workdir: "/tmp/repositories"

# Default is info
# debug, info, error
log_level: "info"

agent:
  # Default prompt template is embed config.
  # prompt_path is a relative path from the execution directory.
  # e.g) config/prompt_template_en.yml
  prompt_path: ""

  # Required
  # LLM model name
  # The recommend model is Claude 3.5 Sonnet
  # If you use AWS Bedrock, set the Model ID
  #   e.g) anthropic.claude-3-5-sonnet-20241022-v2:0
  model: "claude-3-5-sonnet-latest"

  # Maximum steps to run agent
  # The following are defined as 1 step
  # - user to LLM and returned to user from LLM
  # - execution function
  max_steps: 70

  # Skip review agents
  # Default is true
  skip_review_agents: true

  git:
    # git user name
    user_name: "github-actions[bot]"

    # git user email
    user_email: "41898282+github-actions[bot]@users.noreply.github.com"

  # GitHub environment for agent
  github:
    # Don't submit files to GitHub by Pull Request.
    no_submit: false

    # Whether to clone repository to the workdir
    clone_repository: true

    # Required
    # Repositories owner to operate
    owner: ""

  # Allow agent to use function.
  # Belows are the default functions.
  allow_functions:
    - get_pull_request
    # - get_web_page_from_url
    # - get_web_search_result
    - list_files
    - modify_file
    - open_file
    - put_file
    - submit_files
    - search_files
```
