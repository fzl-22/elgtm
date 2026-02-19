//go:build integration

package scm_test

import (
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/fzl-22/elgtm/internal/scm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getGitHubTestConfig(t *testing.T) (token, owner, repo string, prNumber int) {
	token = os.Getenv("GITHUB_TOKEN")
	owner = os.Getenv("GITHUB_TEST_OWNER")
	repo = os.Getenv("GITHUB_TEST_REPO")
	prStr := os.Getenv("GITHUB_TEST_PR")

	if token == "" || owner == "" || repo == "" || prStr == "" {
		t.Skip("Skipping GitHub integration test: Missing required environment variables (GITHUB_TOKEN, GITHUB_TEST_OWNER, GITHUB_TEST_REPO, GITHUB_TEST_PR)")
	}

	prNumber, err := strconv.Atoi(prStr)
	require.NoError(t, err, "GITHUB_TEST_PR must be a valid integer")
	return token, owner, repo, prNumber
}

func TestNewGitHubDriver(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 & time.Second}

	t.Run("Success_InitDriver", func(t *testing.T) {
		token := "fake-github-token"
		driver, err := scm.NewGitHubDriver(httpClient, token)

		assert.NoError(t, err)
		assert.NotNil(t, driver)
	})

	t.Run("Failure_MissingToken", func(t *testing.T) {
		driver, err := scm.NewGitHubDriver(httpClient, "")

		assert.Error(t, err)
		assert.Nil(t, driver)
	})
}

func TestGitHubDriver_Integration(t *testing.T) {
	token, owner, repo, prNumber := getGitHubTestConfig(t)
	t.Logf("Got config: %s %s %s %d", token, owner, repo, prNumber)
}
