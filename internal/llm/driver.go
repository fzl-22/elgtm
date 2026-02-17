package llm

import "context"

type Driver interface {
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
}

type GenerateRequest struct {
	Model            string
	Prompt           string
	ResponseMIMEType string
	Temperature      float32
	MaxTokens        int
}

type GenerateResponse struct {
	Content string
}
