package scm

import (
	"context"
	"fmt"
	"io"
	"net/http"

	scm "github.com/fzl-22/elgtm/internal/scm/types"
	"github.com/google/go-github/v82/github"
)

type GitHubClient struct {
	client     *github.Client
	httpClient *http.Client
	token      string
}

func NewGitHubClient(httpClient *http.Client, token string) *GitHubClient {
	return &GitHubClient{
		client:     github.NewClient(httpClient).WithAuthToken(token),
		httpClient: httpClient,
		token:      token,
	}
}

func (c *GitHubClient) GetPullRequest(ctx context.Context, owner, repo string, number int) (*scm.PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request #%d: %w", number, err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", *pr.DiffURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get diff with status: %d", response.StatusCode)
	}

	diffBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read diff: %w", err)
	}
	rawDiff := string(diffBytes)

	return &scm.PullRequest{
		ID:        *pr.ID,
		Number:    *pr.Number,
		Title:     *pr.Title,
		Body:      *pr.Body,
		URL:       *pr.URL,
		DiffURL:   *pr.DiffURL,
		RawDiff:   rawDiff,
		CreatedAt: pr.CreatedAt.Time,
		UpdatedAt: pr.UpdatedAt.Time,
	}, nil
}
