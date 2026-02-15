package llm

import (
	"context"

	"github.com/fzl-22/elgtm/internal/scm"
)

type LLMClient interface {
	GenerateIssueComment(ctx context.Context, pullRequest scm.PullRequest) (*scm.IssueComment, error)
}
