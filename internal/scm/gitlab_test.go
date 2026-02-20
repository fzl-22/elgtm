package scm_test

import (
	"testing"

	"github.com/fzl-22/elgtm/internal/scm"
	"github.com/stretchr/testify/assert"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestGitLabDriver_NewGitLabDriver(t *testing.T) {
	t.Run("Success_InitDriver", func(t *testing.T) {
		driver, err := scm.NewGitLabDriver("fake-gitlab-token")

		assert.NoError(t, err)
		assert.NotNil(t, driver)
	})

	t.Run("Failure_InitError", func(T *testing.T) {
		badOption := gitlab.WithBaseURL("://invalid-url")

		driver, err := scm.NewGitLabDriver("fake-gitlab-token", badOption)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to initialize gitlab driver")
		assert.Nil(t, driver)
	})
}
