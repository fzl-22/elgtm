package scm

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/google/go-github/v82/github"
)

type GitHubClient struct {
	client     *github.Client
	httpClient *http.Client
	cfg        config.SCM
}

func NewGitHubClient(httpClient *http.Client, cfg config.SCM) *GitHubClient {
	return &GitHubClient{
		client:     github.NewClient(httpClient).WithAuthToken(cfg.Token),
		httpClient: httpClient,
		cfg:        cfg,
	}
}

func (c *GitHubClient) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request #%d: %w", number, err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *pr.DiffURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.cfg.Token))

	res, err := c.httpClient.Do(req)
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

	return &PullRequest{
		ID:        *pr.ID,
		Number:    *pr.Number,
		Title:     *pr.Title,
		Body:      *pr.Body,
		URL:       *pr.URL,
		DiffURL:   *pr.DiffURL,
		RawDiff:   string(diffBytes),
		CreatedAt: pr.CreatedAt.Time,
		UpdatedAt: pr.UpdatedAt.Time,
	}, nil
}
