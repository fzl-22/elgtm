package scm

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"strings"

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
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get merge request: %w", err)
	}

	diffs, _, err := d.client.MergeRequests.ListMergeRequestDiffs(projectPath, int64(req.Number), nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get diff for merge request #%d: %w", req.Number, err)
	}

	var diffBuilder strings.Builder
	for _, diff := range diffs {
		if diffBuilder.Len()+len(diff.Diff) > int(req.MaxDiffSize) {
			diffBuilder.WriteString("\n\n... [DIFF TRUNCATED DUE TO SIZE LIMIT] ...")
			break
		}
		diffBuilder.WriteString(diff.Diff)
	}

	slog.Info("DIFF", "diff", diffBuilder.String())

	parsedMR := &PullRequest{
		ID:        mr.ID,
		Number:    int(mr.IID),
		Title:     mr.Title,
		Body:      mr.Description,
		Author:    mr.Author.Username,
		URL:       mr.WebURL,
		HTMLURL:   mr.WebURL,
		RawDiff:   diffBuilder.String(),
		CreatedAt: *mr.CreatedAt,
		UpdatedAt: *mr.UpdatedAt,
	}
	return &GetPRResponse{
		PR: parsedMR,
	}, nil
}

func (d *GitLabDriver) PostIssueComment(ctx context.Context, req PostIssueCommentRequest) error {
	projectPath := path.Join(req.Owner, req.Repo)
	_, _, err := d.client.Notes.CreateIssueNote(projectPath, int64(req.Number), &gitlab.CreateIssueNoteOptions{
		Body: req.IssueComment.Body,
	})
	if err != nil {
		return fmt.Errorf("failed to create issue note: %w", err)
	}

	return nil
}
