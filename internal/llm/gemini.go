package llm

import (
	"context"
	"fmt"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/scm"
	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	cfg    config.LLM
}

func NewGeminiClient(ctx context.Context, cfg config.LLM) (LLMClient, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("gemini api key is missing")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &GeminiClient{
		client: client,
		cfg:    cfg,
	}, nil
}

func (c *GeminiClient) GenerateIssueComment(ctx context.Context, pullRequest scm.PullRequest) (*scm.IssueComment, error) {
	prompt := fmt.Sprintf("Please review this: %+v", pullRequest)

	temp := c.cfg.Temperature
	generatedContent, err := c.client.Models.GenerateContent(ctx, c.cfg.Model, genai.Text(prompt), &genai.GenerateContentConfig{
		Temperature:      &temp,
		ResponseMIMEType: "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	text := generatedContent.Text()
	return &scm.IssueComment{
		Body: &text,
	}, nil
}
