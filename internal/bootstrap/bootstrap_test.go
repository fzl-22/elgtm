package bootstrap_test

import (
	"context"
	"testing"

	"github.com/fzl-22/elgtm/internal/bootstrap"
	"github.com/fzl-22/elgtm/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestBootstrap_Initialize(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_InitializeEngine", func(t *testing.T) {
		cfg := &config.Config{
			System: config.System{Timeout: 30},
			SCM: config.SCM{
				Platform: config.PlatformGitHub,
				Token:    "fake-token",
			},
			LLM: config.LLM{
				Provider: config.ProviderGemini,
				APIKey:   "fake-api-key",
			},
		}

		engine, err := bootstrap.Initialize(ctx, cfg)

		assert.NoError(t, err)
		assert.NotNil(t, engine)
	})

	t.Run("Failure_UnsupportedSCMPlatform", func(t *testing.T) {
		cfg := &config.Config{
			System: config.System{Timeout: 30},
			SCM: config.SCM{
				Platform: "anything-scm",
				Token:    "fake-token",
			},
			LLM: config.LLM{
				Provider: config.ProviderGemini,
				APIKey:   "fake-api-key",
			},
		}

		engine, err := bootstrap.Initialize(ctx, cfg)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported SCM platform")
		assert.Nil(t, engine)
	})

	t.Run("Failure_SCMDriverInitError", func(t *testing.T) {
		cfg := &config.Config{
			System: config.System{Timeout: 30},
			SCM: config.SCM{
				Platform: config.PlatformGitHub,
				Token:    "",
			},
			LLM: config.LLM{
				Provider: config.ProviderGemini,
				APIKey:   "fake-api-key",
			},
		}

		engine, err := bootstrap.Initialize(ctx, cfg)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to initialize SCM driver")
		assert.Nil(t, engine)
	})

	t.Run("Failure_UnsupportedLLMProvider", func(t *testing.T) {
		cfg := &config.Config{
			System: config.System{Timeout: 30},
			SCM: config.SCM{
				Platform: config.PlatformGitHub,
				Token:    "fake-token",
			},
			LLM: config.LLM{
				Provider: "anything-llm",
				APIKey:   "fake-api-key",
			},
		}

		engine, err := bootstrap.Initialize(ctx, cfg)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported LLM provider")
		assert.Nil(t, engine)
	})

	t.Run("Failure_LLMDriverInitError", func(t *testing.T) {
		t.Setenv("GEMINI_API_KEY", "")

		cfg := &config.Config{
			System: config.System{Timeout: 30},
			SCM: config.SCM{
				Platform: config.PlatformGitHub,
				Token:    "fake-token",
			},
			LLM: config.LLM{
				Provider: config.ProviderGemini,
				APIKey:   "",
			},
		}

		engine, err := bootstrap.Initialize(ctx, cfg)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to initialize LLM driver")
		assert.Nil(t, engine)
	})
}
