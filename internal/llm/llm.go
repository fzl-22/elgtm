package llm

import (
	"context"
)

type Client interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}
