package scm

import (
	"context"
)

type Client interface {
	GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error)
	PostIssueComment(ctx context.Context, owner, repo string, number int, issueComent *IssueComment) error
}
