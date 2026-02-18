package llm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDriver struct {
	mock.Mock
}

func (m *MockDriver) Generate(ctx context.Context, req llm.GenerateRequest) (*llm.GenerateResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*llm.GenerateResponse), args.Error(1)
}

func TestNewClient(t *testing.T) {
	t.Run("Success_InitClient", func(t *testing.T) {
		cfg := config.Config{}

		mockDriver := new(MockDriver)
		client := llm.NewClient(mockDriver, cfg.LLM)

		assert.NotNil(t, client)
		mockDriver.AssertExpectations(t)
	})
}

func TestGenerateContent(t *testing.T) {
	t.Run("Success_SuccessGenerate", func(t *testing.T) {
		cfg := config.Config{}

		mockDriver := new(MockDriver)

		mockDriver.On("Generate", mock.Anything, mock.Anything).
			Return(&llm.GenerateResponse{
				Content: "Looks Good To Me!",
			}, nil)

		client := llm.NewClient(mockDriver, cfg.LLM)

		content, err := client.GenerateContent(context.Background(), "Hi, I am a prompt")

		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Equal(t, content, "Looks Good To Me!")
		mockDriver.AssertExpectations(t)
	})

	t.Run("Failure_FailedToGenerateContent", func(t *testing.T) {
		cfg := config.Config{}

		mockDriver := new(MockDriver)

		mockDriver.On("Generate", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("failed to generate content"))

		client := llm.NewClient(mockDriver, cfg.LLM)

		content, err := client.GenerateContent(context.Background(), "Hi, I am a prompt")

		assert.NotNil(t, client)
		assert.Error(t, err)
		assert.Empty(t, content)
		mockDriver.AssertExpectations(t)
	})
}
