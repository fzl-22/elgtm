package config_test

import (
	"os"
	"testing"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig_ReadsEnvVars(t *testing.T) {
	os.Setenv("LLM_PROVIDER", "claude")
	os.Setenv("LLM_MODEL", "claude-sonnet-4-5-20250929")
	os.Setenv("LLM_API_KEY", "test-api-key")
	os.Setenv("LLM_TEMPERATURE", "0.9") // Default: 0.2
	os.Setenv("LLM_MAX_TOKENS", "1024") // Default: 4096
	os.Setenv("SCM_PLATFORM", "gitlab")
	os.Setenv("SCM_TOKEN", "test-scm-token")
	os.Setenv("SCM_OWNER", "test-owner")
	os.Setenv("SCM_REPO", "test-repo")
	os.Setenv("SCM_PR_NUMBER", "123")
	os.Setenv("SCM_MAX_DIFF_SIZE", "500000")         // Default: 2097152
	os.Setenv("REVIEW_PROMPT_TYPE", "security")      // Default: general
	os.Setenv("REVIEW_PROMPT_DIR", "custom_prompts") // Default: .reviewer
	os.Setenv("REVIEW_LANGUAGE", "id")               // Default: en
	os.Setenv("SYSTEM_LOG_LEVEL", "debug")           // Default: info
	os.Setenv("SYSTEM_TIMEOUT", "60")                // Default: 30
	defer os.Clearenv()

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
}

func TestNewConfig_ReadsEnvVarsWithDefaultValue(t *testing.T) {
	os.Setenv("LLM_PROVIDER", "claude")
	os.Setenv("LLM_MODEL", "claude-sonnet-4-5-20250929")
	os.Setenv("LLM_API_KEY", "test-api-key")
	os.Setenv("SCM_PLATFORM", "gitlab")
	os.Setenv("SCM_TOKEN", "test-scm-token")
	os.Setenv("SCM_OWNER", "test-owner")
	os.Setenv("SCM_REPO", "test-repo")
	os.Setenv("SCM_PR_NUMBER", "123")
	defer os.Clearenv()

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
}
