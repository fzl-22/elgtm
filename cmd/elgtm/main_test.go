package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	setEnv := func(t *testing.T, key, value string) {
		t.Helper()
		t.Setenv(key, value)
	}

	t.Run("Failure_ConfigLoadFailed", func(t *testing.T) {
		os.Clearenv()
		setEnv(t, "SCM_PR_NUMBER", "not-a-number")

		exitCode := run()

		assert.Equal(t, 1, exitCode)
	})

	t.Run("Failure_InitializationFailed", func(t *testing.T) {
		os.Clearenv()
		setEnv(t, "SCM_PLATFORM", "github")
		setEnv(t, "SCM_TOKEN", "")

		exitCode := run()

		assert.Equal(t, 1, exitCode)
	})

	t.Run("Failure_ReviewFailed", func(t *testing.T) {
		os.Clearenv()

		setEnv(t, "SCM_PLATFORM", "github")
		setEnv(t, "SCM_TOKEN", "dummy-token")
		setEnv(t, "LLM_PROVIDER", "gemini")
		setEnv(t, "LLM_API_KEY", "dummy-api-key")
		setEnv(t, "REVIEW_PROMPT_TYPE", "this-prompt-does-not-exist")

		exitCode := run()

		assert.Equal(t, 1, exitCode)
	})
}
