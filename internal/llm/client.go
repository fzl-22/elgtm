package llm

import (
	"context"
	"fmt"

	"github.com/fzl-22/elgtm/internal/config"
)

type client struct {
	driver Driver
	cfg    config.LLM
}

func NewClient(driver Driver, cfg config.LLM) Client {
	return &client{
		driver: driver,
		cfg:    cfg,
	}
}

func (c *client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	req := GenerateRequest{
		Model:       c.cfg.Model,
		Prompt:      prompt,
		Temperature: c.cfg.Temperature,
		MaxTokens:   c.cfg.MaxTokens,
	}

	resp, err := c.driver.Generate(ctx, req)
	if err != nil {
		return "", fmt.Errorf("driver error: %w", err)
	}

	return resp.Content, nil
}
