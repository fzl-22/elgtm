package scm_test

import (
	"context"
	"testing"
	"time"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/scm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDriver struct {
	mock.Mock
}

func (m *MockDriver) GetPullRequest(ctx context.Context, req scm.GetPRRequest) (*scm.GetPRResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*scm.GetPRResponse), args.Error(1)
}

func (m *MockDriver) PostIssueComment(ctx context.Context, req scm.PostIssueCommentRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestNewClient(t *testing.T) {
	t.Run("Success_InitClient", func(t *testing.T) {
		mockDriver := new(MockDriver)

		cfg := config.Config{}
		client := scm.NewClient(mockDriver, cfg.SCM)

		assert.NotNil(t, client)
		mockDriver.AssertExpectations(t)
	})
}

func TestGetPullRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_SuccessGetPullRequest", func(t *testing.T) {
		mockDriver := new(MockDriver)

		fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		expectedPR := &scm.PullRequest{
			ID:        12345,
			URL:       "https://api.github.com/repos/fzl-22/elgtm/pulls/1",
			Title:     "feat: add ai review",
			Number:    1,
			Author:    "fzl-22",
			Body:      "This is a test PR",
			HTMLURL:   "https://github.com/fzl-22/elgtm/pull/1",
			DiffURL:   "https://github.com/fzl-22/elgtm/pull/1.diff",
			RawDiff:   "diff --git a/main.go b/main.go...",
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		}

		mockDriver.On("GetPullRequest", mock.Anything, mock.Anything).
			Return(&scm.GetPRResponse{
				PR: expectedPR,
			}, nil)

		cfg := config.Config{}
		client := scm.NewClient(mockDriver, cfg.SCM)

		pr, err := client.GetPullRequest(ctx, "fzl-22", "elgtm", expectedPR.Number)

		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, expectedPR, pr)
		mockDriver.AssertExpectations(t)
	})

	t.Run("Failure_FailedToGetPullRequest", func(t *testing.T) {
		mockDriver := new(MockDriver)

		mockDriver.On("GetPullRequest", mock.Anything, mock.Anything).
			Return(nil, assert.AnError)

		cfg := config.Config{}
		client := scm.NewClient(mockDriver, cfg.SCM)

		pr, err := client.GetPullRequest(ctx, "fzl-22", "elgtm", 123)

		assert.NotNil(t, client)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get pull request using SCM driver")
		assert.Nil(t, pr)
		mockDriver.AssertExpectations(t)
	})
}

func TestPostIssueComment(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_SuccessPostIssueComment", func(t *testing.T) {
		mockDriver := new(MockDriver)

		mockDriver.On("PostIssueComment", mock.Anything, mock.Anything).
			Return(nil)

		cfg := config.Config{}
		client := scm.NewClient(mockDriver, cfg.SCM)

		commentBody := "Test comment!"
		issueComment := scm.IssueComment{
			Body: &commentBody,
		}
		err := client.PostIssueComment(ctx, "fzl-22", "elgtm", 123, &issueComment)

		assert.NotNil(t, client)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("Failure_FailedToPostIssueComment", func(t *testing.T) {
		mockDriver := new(MockDriver)

		mockDriver.On("PostIssueComment", mock.Anything, mock.Anything).
			Return(assert.AnError)

		cfg := config.Config{}
		client := scm.NewClient(mockDriver, cfg.SCM)

		commentBody := "Test comment!"
		issueComment := scm.IssueComment{
			Body: &commentBody,
		}
		err := client.PostIssueComment(ctx, "fzl-22", "elgtm", 123, &issueComment)

		assert.NotNil(t, client)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to post issue comment using SCM driver")
		mockDriver.AssertExpectations(t)
	})
}
