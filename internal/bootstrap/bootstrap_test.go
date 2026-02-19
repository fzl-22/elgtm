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
}
