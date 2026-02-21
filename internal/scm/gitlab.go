package scm

import (
	"context"
	"fmt"
	"path"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitLabDriver struct {
	client *gitlab.Client
}

func NewGitLabDriver(token string, opts ...gitlab.ClientOptionFunc) (*GitLabDriver, error) {
	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gitlab driver: %w", err)
	}

	return &GitLabDriver{
		client: client,
	}, nil
}

func (d *GitLabDriver) GetPullRequest(ctx context.Context, req GetPRRequest) (*GetPRResponse, error) {
	projectPath := path.Join(req.Owner, req.Repo)

	renderHTML := false
	mr, _, err := d.client.MergeRequests.GetMergeRequest(projectPath, int64(req.Number), &gitlab.GetMergeRequestsOptions{
		RenderHTML: &renderHTML,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	parsedMR := &PullRequest{
		ID:        mr.ID,
		Number:    int(mr.ID),
		Title:     mr.Title,
		Body:      mr.Description,
		Author:    mr.Author.Username,
		URL:       mr.WebURL,
		CreatedAt: *mr.CreatedAt,
		UpdatedAt: *mr.UpdatedAt,
	}
	return &GetPRResponse{
		PR: parsedMR,
	}, nil
}

func (d *GitLabDriver) PostIssueComment(ctx context.Context, req PostIssueCommentRequest) error {
	return nil
}
