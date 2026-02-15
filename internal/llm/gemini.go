package llm

import (
	"context"
	"fmt"

	"github.com/fzl-22/elgtm/internal/config"
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

func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	temp := c.cfg.Temperature
	resp, err := c.client.Models.GenerateContent(ctx, c.cfg.Model, genai.Text(prompt), &genai.GenerateContentConfig{
		Temperature:      &temp,
		ResponseMIMEType: "text/plain",
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return resp.Text(), nil
}
