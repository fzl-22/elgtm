package scm

import (
	"context"
)

type SCMClient interface {
	GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error)
}
