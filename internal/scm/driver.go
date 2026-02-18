package scm

import "context"

type Driver interface {
	GetPullRequest(ctx context.Context, req GetPRRequest) (*GetPRResponse, error)
	PostIssueComment(ctx context.Context, req PostIssueCommentRequest) error
}

type GetPRRequest struct {
	Owner  string
	Repo   string
	Number int
}

type GetPRResponse struct {
	PR *PullRequest
}

type PostIssueCommentRequest struct {
	Owner        string
	Repo         string
	Number       int
	IssueComment *IssueComment
}
