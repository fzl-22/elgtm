package scm

import (
	"context"
)

type SCM interface {
	GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error)
}
