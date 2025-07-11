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
	// ConfigFilePath always has config.yml mounted to
	ConfigFilePath = "/agent/config/config.yml"
	DefaultWorkDir = "/agent/repositories"

	LogDebug = "debug"
	LogInfo  = "info"
	LogError = "error"
)

type Git struct {
	UserName  string `yaml:"user_name"`
	UserEmail string `yaml:"user_email"`
}

type GitHub struct {
	CloneRepository *bool    `yaml:"clone_repository"`
	Owner           string   `yaml:"owner" validate:"required"`
	PRLabels        []string `yaml:"pr_labels"`
}

type Agent struct {
	Model          string   `yaml:"model" validate:"required"`
	MaxSteps       int      `yaml:"max_steps" validate:"gte=0"`
	Git            Git      `yaml:"git"`
	GitHub         GitHub   `yaml:"github"`
	AllowFunctions []string `yaml:"allow_functions"`
}

type Config struct {
	Language string `yaml:"language"`
	WorkDir  string `yaml:"workdir"`
	LogLevel string `yaml:"log_level" validate:"log_level"`
	Agent    Agent  `yaml:"agent" validate:"required"`
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

// LoadInCommand loads the configuration in command mode.
// In command, the config file is mounted to a fixed path.
func LoadInCommand(path string) (Config, error) {
	if path == "" {
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

func Validate(config Config) error {
	validate := validator.New()
	if err := validate.RegisterValidation("log_level", isValidLogLevel); err != nil {
		return err
	}
	if err := validate.Struct(config); err != nil {
		errs := err.(validator.ValidationErrors)
		return fmt.Errorf("validation failed: %w", errs)
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

	if conf.Agent.Git.UserName == "" {
		conf.Agent.Git.UserName = "github-actions[bot]"
	}

	if conf.Agent.Git.UserEmail == "" {
		conf.Agent.Git.UserEmail = "41898282+github-actions[bot]@users.noreply.github.com"
	}

	if conf.Agent.GitHub.CloneRepository == nil {
		clone := true
		conf.Agent.GitHub.CloneRepository = &clone
	}

	return conf
}
