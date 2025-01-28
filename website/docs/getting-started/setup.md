# Setup

## Environment Variables
```shell
# Required for GitHub authentication
GITHUB_TOKEN=your_github_token

# If you use OpenAI models
OPENAI_API_KEY=your_openai_api_key

# If you use Anthropic models
ANTHROPIC_API_KEY=your_anthropic_api_key
```

##  More Configuration
Copy the [default_config.yml](https://github.com/clover0/issue-agent/blob/main/agent/config/default_config.yml) file to your repository root of your repository and rename it to `issue_agent.yml`.
Then, edit the file as needed.


If you prefer to place the `issue_agent.yml` file in a custom path (e.g. `$HOME/your_issue_agent.yml`),
specify the path using the `--config_path` flag.

```shell 

issue-agent \ 
    --config_path /path/to/issue_agent.yml
```
