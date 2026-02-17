package llm

import "context"

type Driver interface {
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
}
