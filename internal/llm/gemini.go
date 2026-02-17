package llm

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiDriver struct {
	client *genai.Client
}

func NewGeminiDriver(ctx context.Context, apiKey string) (*GeminiDriver, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("missing Gemini API Key")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini SDK client: %w", err)
	}

	return &GeminiDriver{client: client}, nil
}

func (d *GeminiDriver) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	sdkConfig := &genai.GenerateContentConfig{
		Temperature:      &req.Temperature,
		MaxOutputTokens:  int32(req.MaxTokens),
		ResponseMIMEType: req.ResponseMIMEType,
	}

	resp, err := d.client.Models.GenerateContent(ctx, req.Model, genai.Text(req.Prompt), sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content from Gemini API: %w", err)
	}

	return &GenerateResponse{
		Content: resp.Text(),
	}, nil
}
