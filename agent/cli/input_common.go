package cli

import (
	"flag"
)

type CommonInput struct {
	Config       string
	AWSProfile   string
	AWSRegion    string
	LogLevel     string
	Language     string
	Model        string
	ReviewAgents int // TODO: remove from common. move to create-pr command.
}

func addCommonFlags(fs *flag.FlagSet, cfg *CommonInput) {
	fs.StringVar(&cfg.Config, "config", "", `Path to the configuration file. 
Default: agent/config/default_config.yml in this project.`)

	fs.StringVar(&cfg.AWSProfile, "aws_profile", "", "AWS profile to use a specific profile from credentials.")

	fs.StringVar(&cfg.AWSRegion, "aws_region", "", `AWS region to use for credentials and Bedrock.
Default(If use aws_profile): aws profile's default session region.`)

	fs.StringVar(&cfg.LogLevel, "log_level", "info", `Log level. If you want to see LLM completions, set it to 'debug'.
Default: info.`)

	fs.StringVar(&cfg.Language, "language", "English", `Language spoken by agent.
Default: English.`)

	fs.StringVar(&cfg.Model, "model", "", "LLM name. For the model name, check the documentation of each LLM provider.")

	fs.IntVar(&cfg.ReviewAgents, "review_agents", 0, `The number of agents to review.
The number of agents to review. A value greater than 0 will review to the created PR.
Default: 0`)
}
