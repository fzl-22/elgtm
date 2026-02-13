package scm

import (
	"context"

	scm "github.com/fzl-22/elgtm/internal/scm/types"
)

type SCM interface {
	GetPullRequest(ctx context.Context, owner, repo string, number int) (*scm.PullRequest, error)
}
