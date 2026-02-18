package scm

import (
	"context"
	"fmt"

	"github.com/fzl-22/elgtm/internal/config"
)

type client struct {
	driver Driver
	cfg    config.SCM
}

func NewClient(driver Driver, cfg config.SCM) Client {
	return &client{
		driver: driver,
		cfg:    cfg,
	}
}

func (c *client) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	req := GetPRRequest{
		Owner:  owner,
		Repo:   repo,
		Number: number,
	}

	resp, err := c.driver.GetPullRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request using SCM driver: %w", err)
	}

	return resp.PR, nil
}

func (c *client) PostIssueComment(ctx context.Context, owner, repo string, number int, issueComment *IssueComment) error {
	req := PostIssueCommentRequest{
		Owner:        owner,
		Repo:         repo,
		Number:       number,
		IssueComment: issueComment,
	}

	err := c.driver.PostIssueComment(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to post issue comment using SCM driver: %w", err)
	}

	return nil
}
