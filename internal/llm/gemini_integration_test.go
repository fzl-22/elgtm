//go:build integration

package llm_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/fzl-22/elgtm/internal/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeminiDriver_NewGeminiDriver(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_InitDriver", func(t *testing.T) {
		apiKey := os.Getenv("GEMINI_API_KEY")
		driver, err := llm.NewGeminiDriver(ctx, apiKey)
		assert.NoError(t, err)
		assert.NotNil(t, driver)
	})

	t.Run("Failure_MissingAPIKey", func(t *testing.T) {
		t.Setenv("GEMINI_API_KEY", "")
		driver, err := llm.NewGeminiDriver(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "api key is required for Google AI backend")
		assert.Nil(t, driver)
	})
}

func TestGeminiDriver_Generate(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GEMINI_API_KEY not set")
	}

	ctx := context.Background()

	t.Run("Success_GenerateContent", func(t *testing.T) {
		driver, err := llm.NewGeminiDriver(ctx, apiKey)
		require.NoError(t, err)

		req := llm.GenerateRequest{
			Model:       "gemini-2.5-flash",
			Prompt:      "Say 'OK'",
			Temperature: 0.1,
			MaxTokens:   128,
		}

		res, err := driver.Generate(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Content)
	})

	t.Run("Failure_InvalidAPIKey", func(t *testing.T) {
		driver, err := llm.NewGeminiDriver(ctx, "invalid-api-key")
		require.NoError(t, err)

		req := llm.GenerateRequest{
			Model:       "gemini-2.5-flash",
			Prompt:      "Say 'OK'",
			Temperature: 0.1,
			MaxTokens:   128,
		}

		res, err := driver.Generate(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API key not valid")
		assert.Empty(t, res)

		t.Logf("Got expected error: %v", err)
	})

	t.Run("Failure_UnknownModel", func(t *testing.T) {
		driver, err := llm.NewGeminiDriver(ctx, apiKey)
		require.NoError(t, err)

		req := llm.GenerateRequest{
			Model:       "gemini-unknown-model",
			Prompt:      "Say 'OK'",
			Temperature: 0.1,
			MaxTokens:   128,
		}

		res, err := driver.Generate(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "models/gemini-unknown-model is not found")
		assert.Empty(t, res)

		t.Logf("Got expected error: %v", err)
	})

	t.Run("Success_ContextTimeout", func(t *testing.T) {
		driver, err := llm.NewGeminiDriver(ctx, "invalid-api-key")
		require.NoError(t, err)

		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		req := llm.GenerateRequest{
			Model:       "gemini-2.5-flash",
			Prompt:      "Say 'OK'",
			Temperature: 0.1,
			MaxTokens:   128,
		}

		res, err := driver.Generate(shortCtx, req)

		assert.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Empty(t, res)

		t.Logf("Got expected error: %v", err)
	})
}
