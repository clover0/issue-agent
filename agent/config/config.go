package config

import (
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yml
var defaultConfig []byte

const (
	ConfigFilePath = "/agent/config/config.yml"
	PromptFilePath = "/agent/config/prompt.yml"
	DefaultWorkDir = "/agent/repositories"

	LogDebug = "debug"
	LogInfo  = "info"
	LogError = "error"
)

type Config struct {
	Language string `yaml:"language"`
	WorkDir  string `yaml:"workdir"`
	LogLevel string `yaml:"log_level" validate:"log_level"`
	Agent    struct {
		PromptPath       string `yaml:"prompt_path"`
		Model            string `yaml:"model"`
		MaxSteps         int    `yaml:"max_steps" validate:"gte=0"`
		SkipReviewAgents *bool  `yaml:"skip_review_agents"`
		Git              struct {
			UserName  string `yaml:"user_name"`
			UserEmail string `yaml:"user_email"`
		} `yaml:"git"`
		GitHub struct {
			NoSubmit        *bool  `yaml:"no_submit"`
			CloneRepository *bool  `yaml:"clone_repository"`
			Owner           string `yaml:"owner" validate:"required"`
		}
		AllowFunctions []string `yaml:"allow_functions"`
	} `yaml:"agent" validate:"required"`
}

func isValidLogLevel(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, level := range []string{LogDebug, LogInfo, LogError} {
		if level == value {
			return true
		}
	}
	return false
}

func LoadDefault(passedConfig bool) (Config, error) {
	if !passedConfig {
		return Load("")
	}

	cf, err := Load(ConfigFilePath)
	if err != nil {
		return cf, err
	}

	return cf, nil
}

func Load(path string) (Config, error) {
	var cnfg Config

	var data []byte
	if path == "" {
		data = defaultConfig
	} else {
		file, err := os.Open(path)
		if err != nil {
			return cnfg, err
		}
		defer file.Close()

		data, err = io.ReadAll(file)
		if err != nil {
			return cnfg, err
		}
	}

	if err := yaml.Unmarshal(data, &cnfg); err != nil {
		return cnfg, err
	}

	cnfg = setDefaults(cnfg)

	return cnfg, nil
}

func ValidateConfig(config Config) error {
	validate := validator.New()
	if err := validate.RegisterValidation("log_level", isValidLogLevel); err != nil {
		return err
	}
	if err := validate.Struct(config); err != nil {
		errs := err.(validator.ValidationErrors)
		return fmt.Errorf("validation failed: %w\n", errs)
	}
	return nil
}

func setDefaults(conf Config) Config {
	if conf.LogLevel == "" {
		conf.LogLevel = LogDebug
	}

	if conf.Language == "" {
		conf.Language = "English"
	}

	if conf.WorkDir == "" {
		conf.WorkDir = DefaultWorkDir
	}

	if conf.Agent.Model == "" {
		conf.Agent.Model = "claude-3-5-sonnet-latest"
	}
	if conf.Agent.MaxSteps == 0 {
		conf.Agent.MaxSteps = 70
	}

	if conf.Agent.Git.UserName == "" {
		conf.Agent.Git.UserName = "github-actions[bot]"
	}
	if conf.Agent.Git.UserEmail == "" {
		conf.Agent.Git.UserEmail = "41898282+github-actions[bot]@users.noreply.github.com"
	}

	if conf.Agent.GitHub.NoSubmit == nil {
		noSubmit := false
		conf.Agent.GitHub.NoSubmit = &noSubmit
	}
	if conf.Agent.GitHub.CloneRepository == nil {
		clone := true
		conf.Agent.GitHub.CloneRepository = &clone
	}

	if conf.Agent.SkipReviewAgents == nil {
		skip := true
		conf.Agent.SkipReviewAgents = &skip
	}

	// TODO: default value
	if len(conf.Agent.AllowFunctions) == 0 {
		conf.Agent.AllowFunctions = []string{
			"submit_files",
			"get_pull_request_diff",
			// "get_web_page_from_url",
			// "get_web_search_result",
			"list_files",
			"modify_file",
			"open_file",
			"put_file",
			"submit_files",
			"search_files",
		}
	}

	return conf
}
