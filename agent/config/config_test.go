package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/clover0/issue-agent/config"
	"github.com/clover0/issue-agent/test/assert"
)

func TestIsValidLogLevel(t *testing.T) {
	t.Parallel()

	validate := validator.New()
	if err := validate.RegisterValidation("log_level", config.IsValidLogLevel); err != nil {
		t.Fatalf("failed to register validation: %v", err)
	}

	tests := map[string]struct {
		logLevel string
		valid    bool
	}{
		"debug level": {
			logLevel: config.LogDebug,
			valid:    true,
		},
		"info level": {
			logLevel: config.LogInfo,
			valid:    true,
		},
		"error level": {
			logLevel: config.LogError,
			valid:    true,
		},
		"invalid level": {
			logLevel: "warning",
			valid:    false,
		},
		"empty level": {
			logLevel: "",
			valid:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			type testStruct struct {
				LogLevel string `validate:"log_level"`
			}

			ts := testStruct{LogLevel: tt.logLevel}

			err := validate.Struct(ts)

			if tt.valid {
				assert.Nil(t, err)
				return
			}
			assert.HasError(t, err)
		})
	}
}

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("default config", func(t *testing.T) {
		t.Parallel()

		cfg, err := config.Load("")
		assert.Nil(t, err)
		assert.Equal(t, cfg.LogLevel, config.LogInfo)
		assert.Equal(t, cfg.Language, "English")
		assert.Equal(t, cfg.WorkDir, "/tmp/repositories")
		assert.Equal(t, cfg.Agent.MaxSteps, 70)
		assert.Equal(t, cfg.Agent.Git.UserName, "github-actions[bot]")
		assert.Equal(t, cfg.Agent.Git.UserEmail, "41898282+github-actions[bot]@users.noreply.github.com")
		assert.Equal(t, *cfg.Agent.GitHub.NoSubmit, false)
		assert.Equal(t, *cfg.Agent.GitHub.CloneRepository, true)
		if len(cfg.Agent.AllowFunctions) == 0 {
			t.Errorf("wanted AllowFunctions to have elements, but it was empty")
		}
	})

	t.Run("load from file", func(t *testing.T) {
		t.Parallel()

		// Create a temporary config file
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yml")
		configContent := `
language: Japanese
workdir: /custom/workdir
log_level: info
agent:
  model: gpt-4
  max_steps: 50
  git:
    user_name: test-user
    user_email: test@example.com
  github:
    owner: test-owner
    pr_labels:
      - test-label
`
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		assert.Nil(t, err)

		cfg, err := config.Load(configPath)

		assert.Nil(t, err)
		assert.Equal(t, cfg.Language, "Japanese")
		assert.Equal(t, cfg.WorkDir, "/custom/workdir")
		assert.Equal(t, cfg.LogLevel, config.LogInfo)
		assert.Equal(t, cfg.Agent.Model, "gpt-4")
		assert.Equal(t, cfg.Agent.MaxSteps, 50)
		assert.Equal(t, cfg.Agent.Git.UserName, "test-user")
		assert.Equal(t, cfg.Agent.Git.UserEmail, "test@example.com")
		assert.Equal(t, cfg.Agent.GitHub.Owner, "test-owner")
		assert.Equal(t, len(cfg.Agent.GitHub.PRLabels), 1)
		assert.Equal(t, cfg.Agent.GitHub.PRLabels[0], "test-label")
	})

	t.Run("non-existent file", func(t *testing.T) {
		t.Parallel()

		_, err := config.Load("/non/existent/path")
		assert.HasError(t, err)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		t.Parallel()

		// Create a temporary config file with invalid YAML
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "invalid.yml")
		configContent := `
language: Japanese
workdir: /custom/workdir
log_level: info
agent:
  model: gpt-4
  max_steps: 50
  invalid yaml
`
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		assert.Nil(t, err)

		_, err = config.Load(configPath)

		assert.HasError(t, err)
	})
}

func TestLoadInCommand(t *testing.T) {
	t.Parallel()

	t.Run("empty path", func(t *testing.T) {
		t.Parallel()

		cfg, err := config.LoadInCommand("")
		assert.Nil(t, err)
		assert.Equal(t, cfg.LogLevel, config.LogInfo)
	})

	// Note: Testing with a non-empty path would require mocking the ConfigFilePath,
	// which might be beyond the scope of this test implementation.
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid config", func(t *testing.T) {
		t.Parallel()

		cfg := config.Config{
			LogLevel: config.LogDebug,
			Agent: config.AgentConfig{
				Model: "gpt-4",
				GitHub: config.GitHubConfig{
					Owner: "test-owner",
				},
			},
		}

		err := config.ValidateConfig(cfg)

		assert.Nil(t, err)
	})

	t.Run("invalid log level", func(t *testing.T) {
		t.Parallel()

		cfg := config.Config{
			LogLevel: "warning", // Invalid log level
			Agent: config.AgentConfig{
				Model: "gpt-4",
				GitHub: config.GitHubConfig{
					Owner: "test-owner",
				},
			},
		}

		err := config.ValidateConfig(cfg)

		assert.HasError(t, err)
	})

	t.Run("missing model field", func(t *testing.T) {
		t.Parallel()

		cfg := config.Config{
			LogLevel: config.LogDebug,
			Agent: config.AgentConfig{
				GitHub: config.GitHubConfig{
					Owner: "test-owner",
				},
			},
		}

		err := config.ValidateConfig(cfg)

		assert.HasError(t, err)
	})

	t.Run("missing owner field", func(t *testing.T) {
		t.Parallel()

		cfg := config.Config{
			LogLevel: config.LogDebug,
			Agent: config.AgentConfig{
				Model:  "gpt-4",
				GitHub: config.GitHubConfig{
					// Owner is missing
				},
			},
		}

		err := config.ValidateConfig(cfg)

		assert.HasError(t, err)
	})
}

func TestSetDefaults(t *testing.T) {
	t.Parallel()

	t.Run("empty config", func(t *testing.T) {
		t.Parallel()

		cfg := config.Config{}
		cfg = config.SetDefaults(cfg)

		assert.Equal(t, cfg.LogLevel, config.LogDebug)
		assert.Equal(t, cfg.Language, "English")
		assert.Equal(t, cfg.WorkDir, config.DefaultWorkDir)
		assert.Equal(t, cfg.Agent.Git.UserName, "github-actions[bot]")
		assert.Equal(t, cfg.Agent.Git.UserEmail, "41898282+github-actions[bot]@users.noreply.github.com")
		assert.Equal(t, *cfg.Agent.GitHub.NoSubmit, false)
		assert.Equal(t, *cfg.Agent.GitHub.CloneRepository, true)
	})

	t.Run("preserve existing values", func(t *testing.T) {
		t.Parallel()

		noSubmit := true
		clone := false
		cfg := config.Config{
			LogLevel: config.LogInfo,
			Language: "Japanese",
			WorkDir:  "/custom/workdir",
			Agent: config.AgentConfig{
				MaxSteps: 50,
				Git: config.GitConfig{
					UserName:  "test-user",
					UserEmail: "test@example.com",
				},
				GitHub: config.GitHubConfig{
					NoSubmit:        &noSubmit,
					CloneRepository: &clone,
					Owner:           "test-owner",
				},
				AllowFunctions: []string{"custom-function"},
			},
		}

		cfg = config.SetDefaults(cfg)

		assert.Equal(t, cfg.LogLevel, config.LogInfo)
		assert.Equal(t, cfg.Language, "Japanese")
		assert.Equal(t, cfg.WorkDir, "/custom/workdir")
		assert.Equal(t, cfg.Agent.MaxSteps, 50)
		assert.Equal(t, cfg.Agent.Git.UserName, "test-user")
		assert.Equal(t, cfg.Agent.Git.UserEmail, "test@example.com")
		assert.Equal(t, *cfg.Agent.GitHub.NoSubmit, true)
		assert.Equal(t, *cfg.Agent.GitHub.CloneRepository, false)
		assert.Equal(t, len(cfg.Agent.AllowFunctions), 1)
		assert.Equal(t, cfg.Agent.AllowFunctions[0], "custom-function")
	})
}
