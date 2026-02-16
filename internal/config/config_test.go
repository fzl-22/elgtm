package config_test

import (
	"os"
	"testing"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	setEnv := func(t *testing.T, key, value string) {
		t.Helper()
		os.Setenv(key, value)
	}

	t.Run("Success_OverrideAllDefaults", func(t *testing.T) {
		os.Clearenv()
		defer os.Clearenv()

		setEnv(t, "LLM_PROVIDER", "claude")
		setEnv(t, "LLM_MODEL", "claude-sonnet-4-5-20250929")
		setEnv(t, "LLM_API_KEY", "test-api-key")
		setEnv(t, "LLM_TEMPERATURE", "0.9") // Default: 0.2
		setEnv(t, "LLM_MAX_TOKENS", "1024") // Default: 4096
		setEnv(t, "SCM_PLATFORM", "gitlab")
		setEnv(t, "SCM_TOKEN", "test-scm-token")
		setEnv(t, "SCM_OWNER", "test-owner")
		setEnv(t, "SCM_REPO", "test-repo")
		setEnv(t, "SCM_PR_NUMBER", "123")
		setEnv(t, "SCM_MAX_DIFF_SIZE", "500000")         // Default: 2097152
		setEnv(t, "REVIEW_PROMPT_TYPE", "security")      // Default: general
		setEnv(t, "REVIEW_PROMPT_DIR", "custom_prompts") // Default: .reviewer
		setEnv(t, "REVIEW_LANGUAGE", "id")               // Default: en
		setEnv(t, "SYSTEM_LOG_LEVEL", "debug")           // Default: info
		setEnv(t, "SYSTEM_TIMEOUT", "60")                // Default: 30

		cfg, err := config.NewConfig()

		assert.NoError(t, err)
		assert.Equal(t, config.LLMProvider("claude"), cfg.LLM.Provider)
		assert.Equal(t, "claude-sonnet-4-5-20250929", cfg.LLM.Model)
		assert.Equal(t, "test-api-key", cfg.LLM.APIKey)
		assert.Equal(t, float32(0.9), cfg.LLM.Temperature)
		assert.Equal(t, 1024, cfg.LLM.MaxTokens)
		assert.Equal(t, config.SCMPlatform("gitlab"), cfg.SCM.Platform)
		assert.Equal(t, "test-scm-token", cfg.SCM.Token)
		assert.Equal(t, "test-owner", cfg.SCM.Owner)
		assert.Equal(t, "test-repo", cfg.SCM.Repo)
		assert.Equal(t, 123, cfg.SCM.PRNumber)
		assert.Equal(t, int64(500000), cfg.SCM.MaxDiffSize)
		assert.Equal(t, "security", cfg.Review.PromptType)
		assert.Equal(t, "custom_prompts", cfg.Review.PromptDir)
		assert.Equal(t, "id", cfg.Review.Language)
		assert.Equal(t, "debug", cfg.System.LogLevel)
		assert.Equal(t, 60, cfg.System.Timeout)
	})

	t.Run("Success_UseDefaultValues", func(t *testing.T) {
		os.Clearenv()
		defer os.Clearenv()

		setEnv(t, "LLM_PROVIDER", "claude")
		setEnv(t, "LLM_MODEL", "claude-sonnet-4-5-20250929")
		setEnv(t, "LLM_API_KEY", "test-api-key")
		setEnv(t, "SCM_PLATFORM", "gitlab")
		setEnv(t, "SCM_TOKEN", "test-scm-token")
		setEnv(t, "SCM_OWNER", "test-owner")
		setEnv(t, "SCM_REPO", "test-repo")
		setEnv(t, "SCM_PR_NUMBER", "123")

		cfg, err := config.NewConfig()

		assert.NoError(t, err)
		assert.Equal(t, config.LLMProvider("claude"), cfg.LLM.Provider)
		assert.Equal(t, "claude-sonnet-4-5-20250929", cfg.LLM.Model)
		assert.Equal(t, "test-api-key", cfg.LLM.APIKey)
		assert.Equal(t, float32(0.2), cfg.LLM.Temperature)
		assert.Equal(t, 4096, cfg.LLM.MaxTokens)
		assert.Equal(t, config.SCMPlatform("gitlab"), cfg.SCM.Platform)
		assert.Equal(t, "test-scm-token", cfg.SCM.Token)
		assert.Equal(t, "test-owner", cfg.SCM.Owner)
		assert.Equal(t, "test-repo", cfg.SCM.Repo)
		assert.Equal(t, 123, cfg.SCM.PRNumber)
		assert.Equal(t, int64(2097152), cfg.SCM.MaxDiffSize)
		assert.Equal(t, "general", cfg.Review.PromptType)
		assert.Equal(t, ".reviewer", cfg.Review.PromptDir)
		assert.Equal(t, "en", cfg.Review.Language)
		assert.Equal(t, "info", cfg.System.LogLevel)
		assert.Equal(t, 300, cfg.System.Timeout)
	})
}
