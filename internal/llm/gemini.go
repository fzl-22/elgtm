package llm

import (
	"context"

	"google.golang.org/genai"
)

type GeminiDriver struct {
	client *genai.Client
}

func NewGeminiDriver(ctx context.Context, apiKey string) (*GeminiDriver, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiDriver{client: client}, nil
}

func (d *GeminiDriver) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	sdkConfig := &genai.GenerateContentConfig{
		Temperature:      &req.Temperature,
		MaxOutputTokens:  int32(req.MaxTokens),
		ResponseMIMEType: "text/plain",
	}

	resp, err := d.client.Models.GenerateContent(ctx, req.Model, genai.Text(req.Prompt), sdkConfig)
	if err != nil {
		return nil, err
	}

	return &GenerateResponse{
		Content: resp.Text(),
	}, nil
}
