package scm

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/google/go-github/v82/github"
)

type GitHubDriver struct {
	client     *github.Client
	httpClient *http.Client
	cfg        config.SCM
}

func NewGitHubDriver(httpClient *http.Client, cfg config.SCM) (*GitHubDriver, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("github token is missing")
	}

	return &GitHubDriver{
		client:     github.NewClient(httpClient).WithAuthToken(cfg.Token),
		httpClient: httpClient,
		cfg:        cfg,
	}, nil
}

func (c *GitHubDriver) GetPullRequest(ctx context.Context, req GetPRRequest) (*GetPRResponse, error) {
	pr, _, err := c.client.PullRequests.Get(ctx, req.Owner, req.Repo, req.Number)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request #%d: %w", req.Number, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *pr.DiffURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("token %s", c.cfg.Token))

	res, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get diff with status: %d", res.StatusCode)
	}

	limitedReader := io.LimitReader(res.Body, c.cfg.MaxDiffSize)

	diffBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read diff: %w", err)
	}

	if int64(len(diffBytes)) == c.cfg.MaxDiffSize {
		truncationMessage := "\n\n... [DIFF TRUNCATED DUE TO SIZE LIMIT] ..."
		diffBytes = append(diffBytes, []byte(truncationMessage)...)
	}

	parsedPR := &PullRequest{
		ID:        pr.GetID(),
		Number:    pr.GetNumber(),
		Title:     pr.GetTitle(),
		Body:      pr.GetBody(),
		Author:    pr.GetUser().GetLogin(),
		URL:       pr.GetURL(),
		HTMLURL:   pr.GetHTMLURL(),
		DiffURL:   pr.GetDiffURL(),
		RawDiff:   string(diffBytes),
		CreatedAt: pr.GetCreatedAt().Time,
		UpdatedAt: pr.GetUpdatedAt().Time,
	}

	return &GetPRResponse{
		PR: parsedPR,
	}, nil
}

func (c *GitHubDriver) PostIssueComment(ctx context.Context, req PostIssueCommentRequest) error {
	issueComment := github.IssueComment{
		Body: req.IssueComment.Body,
	}
	_, _, err := c.client.Issues.CreateComment(ctx, req.Owner, req.Repo, req.Number, &issueComment)
	if err != nil {
		return fmt.Errorf("failed to post issue comment: %w", err)
	}

	return nil
}
