package llm

import (
	"context"
)

type LLMClient interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}
