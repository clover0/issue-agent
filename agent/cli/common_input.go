package cli

import (
	"flag"
)

type CommonInput struct {
	Config string
}

func addCommonFlags(fs *flag.FlagSet, cfg *CommonInput) {
	fs.StringVar(&cfg.Config, "config", "", "Path to the configuration file. Default is `agent/config/default_config.yml in this project`")
}
