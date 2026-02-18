package reviewer_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/reviewer"
	"github.com/fzl-22/elgtm/internal/scm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSCMClient struct {
	mock.Mock
}

func (m *MockSCMClient) GetPullRequest(ctx context.Context, owner, repo string, number int) (*scm.PullRequest, error) {
	args := m.Called(ctx, owner, repo, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*scm.PullRequest), args.Error(1)
}

func (m *MockSCMClient) PostIssueComment(ctx context.Context, owner, repo string, number int, comment *scm.IssueComment) error {
	args := m.Called(ctx, owner, repo, number, comment)
	return args.Error(0)
}

type MockLLMClient struct {
	mock.Mock
}

func (m *MockLLMClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func TestNewEngine(t *testing.T) {
	t.Run("Success_InitEngine", func(t *testing.T) {
		cfg := config.Config{}

		mockSCMClient := new(MockSCMClient)
		mockLLMClient := new(MockLLMClient)

		engine := reviewer.NewEngine(cfg, mockSCMClient, mockLLMClient)

		assert.NotNil(t, engine)
		mockSCMClient.AssertExpectations(t)
		mockLLMClient.AssertExpectations(t)
	})
}

func TestRun(t *testing.T) {
	createFile := func(t *testing.T, path string, content []byte) {
		t.Helper()
		err := os.WriteFile(path, content, 0644)
		assert.NoError(t, err)
	}

	t.Run("Success_SuccessRunEngine", func(t *testing.T) {
		tempDir := t.TempDir()
		promptPath := filepath.Join(tempDir, "general.md")
		createFile(t, promptPath, []byte("Hello! This is PR {{ .Number }}"))

		cfg := config.Config{
			SCM: config.SCM{
				Owner:    "owner",
				Repo:     "repo",
				PRNumber: 123,
			},
			Review: config.Review{
				PromptType: "general",
				PromptDir:  tempDir,
			},
			System: config.System{
				LogLevel: "debug",
			},
		}

		mockSCMClient := new(MockSCMClient)
		mockLLMClient := new(MockLLMClient)

		mockSCMClient.On("GetPullRequest", mock.Anything, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber).
			Return(&scm.PullRequest{
				Number: 123,
			}, nil)

		mockLLMClient.On("GenerateContent", mock.Anything, mock.Anything).
			Return("Looks Good To Me!", nil)

		mockSCMClient.On("PostIssueComment", mock.Anything, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber, mock.MatchedBy(func(c *scm.IssueComment) bool {
			return *c.Body == "Looks Good To Me!"
		})).Return(nil)

		engine := reviewer.NewEngine(cfg, mockSCMClient, mockLLMClient)

		err := engine.Run(context.Background())

		assert.NotNil(t, engine)
		assert.NoError(t, err)
		mockSCMClient.AssertExpectations(t)
		mockLLMClient.AssertExpectations(t)
	})

	t.Run("Failure_PromptResolutionFailed", func(t *testing.T) {
		cfg := config.Config{
			Review: config.Review{
				PromptType: "missing_prompt",
				PromptDir:  "prompt_dir",
			},
		}

		mockSCMClient := new(MockSCMClient)
		mockLLMClient := new(MockLLMClient)

		engine := reviewer.NewEngine(cfg, mockSCMClient, mockLLMClient)

		err := engine.Run(context.Background())

		assert.NotNil(t, engine)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "prompt resolution failed")
		mockSCMClient.AssertExpectations(t)
		mockLLMClient.AssertExpectations(t)
	})

	t.Run("Failure_FailedToLoadPromptFile", func(t *testing.T) {
		tempDir := t.TempDir()
		promptPath := filepath.Join(tempDir, "general.md")
		createFile(t, promptPath, []byte("content"))

		err := os.Chmod(promptPath, 0000)
		assert.NoError(t, err)

		cfg := config.Config{
			Review: config.Review{
				PromptType: "general",
				PromptDir:  tempDir,
			},
		}

		mockSCMClient := new(MockSCMClient)
		mockLLMClient := new(MockLLMClient)

		engine := reviewer.NewEngine(cfg, mockSCMClient, mockLLMClient)

		err = engine.Run(context.Background())

		assert.NotNil(t, engine)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load prompt file")
		mockSCMClient.AssertExpectations(t)
		mockLLMClient.AssertExpectations(t)
	})

	t.Run("Failure_FailedToGetPullRequest", func(t *testing.T) {
		tempDir := t.TempDir()
		promptPath := filepath.Join(tempDir, "general.md")
		createFile(t, promptPath, []byte("content"))

		cfg := config.Config{
			SCM: config.SCM{
				Owner:    "owner",
				Repo:     "repo",
				PRNumber: 1,
			},
			Review: config.Review{
				PromptType: "general",
				PromptDir:  tempDir,
			},
		}

		mockSCMClient := new(MockSCMClient)
		mockLLMClient := new(MockLLMClient)

		mockSCMClient.On("GetPullRequest", mock.Anything, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber).
			Return(nil, fmt.Errorf("failed to get pull request"))

		engine := reviewer.NewEngine(cfg, mockSCMClient, mockLLMClient)

		err := engine.Run(context.Background())

		assert.NotNil(t, engine)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get pull request")
		mockSCMClient.AssertExpectations(t)
		mockLLMClient.AssertExpectations(t)
	})

	t.Run("Failure_PromptGenerationFailed", func(t *testing.T) {
		tempDir := t.TempDir()
		promptPath := filepath.Join(tempDir, "general.md")
		createFile(t, promptPath, []byte("Hello! This is PR {{ .Number "))

		cfg := config.Config{
			SCM: config.SCM{
				Owner:    "owner",
				Repo:     "repo",
				PRNumber: 123,
			},
			Review: config.Review{
				PromptType: "general",
				PromptDir:  tempDir,
			},
		}

		mockSCMClient := new(MockSCMClient)
		mockLLMClient := new(MockLLMClient)

		mockSCMClient.On("GetPullRequest", mock.Anything, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber).
			Return(&scm.PullRequest{
				Number: 123,
			}, nil)

		engine := reviewer.NewEngine(cfg, mockSCMClient, mockLLMClient)

		err := engine.Run(context.Background())

		assert.NotNil(t, engine)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "prompt generation failed")
		mockSCMClient.AssertExpectations(t)
		mockLLMClient.AssertExpectations(t)
	})

	t.Run("Failure_FailedToGenerateReview", func(t *testing.T) {
		tempDir := t.TempDir()
		promptPath := filepath.Join(tempDir, "general.md")
		createFile(t, promptPath, []byte("Hello! This is PR {{ .Number }}"))

		cfg := config.Config{
			SCM: config.SCM{
				Owner:    "owner",
				Repo:     "repo",
				PRNumber: 123,
			},
			Review: config.Review{
				PromptType: "general",
				PromptDir:  tempDir,
			},
		}

		mockSCMClient := new(MockSCMClient)
		mockLLMClient := new(MockLLMClient)

		mockSCMClient.On("GetPullRequest", mock.Anything, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber).
			Return(&scm.PullRequest{
				Number: 123,
			}, nil)

		mockLLMClient.On("GenerateContent", mock.Anything, mock.AnythingOfType("string")).
			Return("", fmt.Errorf("failed to generate content"))

		engine := reviewer.NewEngine(cfg, mockSCMClient, mockLLMClient)

		err := engine.Run(context.Background())

		assert.NotNil(t, engine)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to generate review")
		mockSCMClient.AssertExpectations(t)
		mockLLMClient.AssertExpectations(t)
	})

	t.Run("Failure_FailedToPostIssueComment", func(t *testing.T) {
		tempDir := t.TempDir()
		promptPath := filepath.Join(tempDir, "general.md")
		createFile(t, promptPath, []byte("Hello! This is PR {{ .Number }}"))

		cfg := config.Config{
			SCM: config.SCM{
				Owner:    "owner",
				Repo:     "repo",
				PRNumber: 123,
			},
			Review: config.Review{
				PromptType: "general",
				PromptDir:  tempDir,
			},
		}

		mockSCMClient := new(MockSCMClient)
		mockLLMClient := new(MockLLMClient)

		mockSCMClient.On("GetPullRequest", mock.Anything, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber).
			Return(&scm.PullRequest{
				Number: 123,
			}, nil)

		mockLLMClient.On("GenerateContent", mock.Anything, mock.AnythingOfType("string")).
			Return("Looks Good To Me!", nil)

		mockSCMClient.On("PostIssueComment", mock.Anything, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber, mock.Anything).
			Return(fmt.Errorf("failed to post issue comment"))

		engine := reviewer.NewEngine(cfg, mockSCMClient, mockLLMClient)

		err := engine.Run(context.Background())

		assert.NotNil(t, engine)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to post issue comment")
		mockSCMClient.AssertExpectations(t)
		mockLLMClient.AssertExpectations(t)
	})
}

func TestResolvePromptPath(t *testing.T) {
	createFile := func(t *testing.T, dir, name string) {
		t.Helper()
		path := filepath.Join(dir, name)
		err := os.WriteFile(path, []byte("content"), 0644)
		assert.NoError(t, err)
	}

	t.Run("Success_FoundInUserDir", func(t *testing.T) {
		userDir := t.TempDir()
		createFile(t, userDir, "general.md")

		e := &reviewer.Engine{}

		path, err := e.ResolvePromptPath(userDir, "general")

		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(userDir, "general.md"), path)
	})

	t.Run("Success_FoundInSystemDir", func(t *testing.T) {
		userDir := t.TempDir()
		systemDir := t.TempDir()
		createFile(t, systemDir, "general.md")

		t.Setenv("PROMPT_DEFAULTS", systemDir)

		e := &reviewer.Engine{}

		path, err := e.ResolvePromptPath(userDir, "general")

		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(systemDir, "general.md"), path)
	})

	t.Run("Failure_NotFoundAnywhere", func(t *testing.T) {
		userDir := t.TempDir()

		t.Setenv("PROMPT_DEFAULTS", userDir)

		e := &reviewer.Engine{}

		path, err := e.ResolvePromptPath(userDir, "missing_prompt")

		assert.Error(t, err)
		assert.Empty(t, path)
		assert.Contains(t, err.Error(), "not found in local")
	})
}
