package scm_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/fzl-22/elgtm/internal/scm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func (e *errorReader) Close() error {
	return nil
}

func TestGitHubDriver_NewGitHubDriver(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}

	t.Run("Success_InitDriver", func(t *testing.T) {
		token := "fake-github-token"
		driver, err := scm.NewGitHubDriver(httpClient, token)

		assert.NoError(t, err)
		assert.NotNil(t, driver)
	})

	t.Run("Failure_MissingToken", func(t *testing.T) {
		driver, err := scm.NewGitHubDriver(httpClient, "")

		assert.Error(t, err)
		assert.Nil(t, driver)
	})
}

func TestGitHubDriver_GetPullRequest(t *testing.T) {
	ctx := context.Background()
	req := scm.GetPRRequest{
		Owner:       "owner",
		Repo:        "repo",
		Number:      1,
		Token:       "fake-token",
		MaxDiffSize: 1024,
	}

	t.Run("Success_GetPullRequest", func(t *testing.T) {
		fakeDiffContent := "diff --git a/main.go b/main.go\n+ fmt.Println(\"Hello ELGTM\")"
		fakePRJson := `{
			"id": 12345,
			"number": 1,
			"title": "feat: add ai review",
			"body": "This is a test PR",
			"user": {"login": "fzl-22"},
			"url": "https://api.github.com/repos/owner/repo/pulls/1",
			"html_url": "https://github.com/owner/repo/pull/1",
			"diff_url": "https://github.com/fake/diff",
			"created_at": "2024-01-01T12:00:00Z",
			"updated_at": "2024-01-01T12:00:00Z"
		}`

		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.Path, "/pulls/") {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(fakePRJson)),
						Header:     make(http.Header),
					}, nil
				}

				if req.URL.String() == "https://github.com/fake/diff" {
					assert.Equal(t, "token fake-token", req.Header.Get("Authorization"))

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(fakeDiffContent)),
						Header:     make(http.Header),
					}, nil
				}

				return nil, errors.New("unexpected HTTP call: " + req.URL.String())
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "fake-token")
		require.NoError(t, err)

		res, err := driver.GetPullRequest(ctx, req)

		assert.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.PR)

		assert.Equal(t, int64(12345), res.PR.ID)
		assert.Equal(t, 1, res.PR.Number)
		assert.Equal(t, "feat: add ai review", res.PR.Title)
		assert.Equal(t, "This is a test PR", res.PR.Body)
		assert.Equal(t, "fzl-22", res.PR.Author)
		assert.Equal(t, "https://github.com/fake/diff", res.PR.DiffURL)

		assert.Equal(t, fakeDiffContent, res.PR.RawDiff)
	})

	t.Run("Success_GetPullRequestWithTruncation", func(t *testing.T) {
		fakeDiffContent := "diff --git a/main.go b/main.go\n+ fmt.Println(\"Hello ELGTM\")"
		fakePRJson := `{
			"id": 12345,
			"number": 1,
			"title": "feat: add ai review",
			"body": "This is a test PR",
			"user": {"login": "fzl-22"},
			"url": "https://api.github.com/repos/owner/repo/pulls/1",
			"html_url": "https://github.com/owner/repo/pull/1",
			"diff_url": "https://github.com/fake/diff",
			"created_at": "2024-01-01T12:00:00Z",
			"updated_at": "2024-01-01T12:00:00Z"
		}`

		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.Path, "/pulls/") {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(fakePRJson)),
						Header:     make(http.Header),
					}, nil
				}

				if req.URL.String() == "https://github.com/fake/diff" {
					assert.Equal(t, "token fake-token", req.Header.Get("Authorization"))

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(fakeDiffContent)),
						Header:     make(http.Header),
					}, nil
				}

				return nil, errors.New("unexpected HTTP call: " + req.URL.String())
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "fake-token")
		require.NoError(t, err)

		// make MaxDiffSize small to trigger truncation
		req.MaxDiffSize = 12
		res, err := driver.GetPullRequest(ctx, req)

		assert.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.PR)

		assert.Equal(t, int64(12345), res.PR.ID)
		assert.Equal(t, 1, res.PR.Number)
		assert.Equal(t, "feat: add ai review", res.PR.Title)
		assert.Equal(t, "This is a test PR", res.PR.Body)
		assert.Equal(t, "fzl-22", res.PR.Author)
		assert.Equal(t, "https://github.com/fake/diff", res.PR.DiffURL)

		assert.Equal(t, fakeDiffContent[:req.MaxDiffSize]+"\n\n... [DIFF TRUNCATED DUE TO SIZE LIMIT] ...", res.PR.RawDiff)
	})

	t.Run("Failure_FailedToGetPullRequest", func(t *testing.T) {

		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"message": "Not Found"}`)),
					Header:     make(http.Header),
				}, nil
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "token")
		require.NoError(t, err)

		res, err := driver.GetPullRequest(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "failed to get pull request")
	})

	t.Run("Failure_FailedToCreateRequest", func(t *testing.T) {
		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				fakePR := `{"id": 1, "number": 1, "diff_url": "://invalid-url"}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(fakePR)),
					Header:     make(http.Header),
				}, nil
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "token")
		require.NoError(t, err)

		res, err := driver.GetPullRequest(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "failed to create request")
	})

	t.Run("Failure_FailedToGetDiff", func(t *testing.T) {
		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.Path, "/pulls/") {
					fakePR := `{"id": 1, "number": 1, "diff_url": "https://github.com/fake/diff"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(fakePR)),
						Header:     make(http.Header),
					}, nil
				}

				return nil, errors.New("simulated network timeout")
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "token")
		require.NoError(t, err)

		res, err := driver.GetPullRequest(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "failed to get diff")
	})

	t.Run("Failure_FailedToGetDiff", func(t *testing.T) {
		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.Path, "/pulls/") {
					fakePR := `{"id": 1, "number": 1, "diff_url": "https://github.com/fake/diff"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(fakePR)),
						Header:     make(http.Header),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`Internal Server Error`)),
					Header:     make(http.Header),
				}, nil
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "token")
		require.NoError(t, err)

		res, err := driver.GetPullRequest(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "failed to get diff with status: 500")
	})

	t.Run("Failure_FailedToReadDiff", func(t *testing.T) {
		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.Path, "/pulls/") {
					fakePR := `{"id": 1, "number": 1, "diff_url": "https://github.com/fake/diff"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(fakePR)),
						Header:     make(http.Header),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       &errorReader{},
					Header:     make(http.Header),
				}, nil
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "token")
		require.NoError(t, err)

		res, err := driver.GetPullRequest(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "failed to read diff")
	})
}

func TestGitHubDriver_PostIssueComment(t *testing.T) {
	ctx := context.Background()
	bodyStr := "This is an AI review comment"

	req := scm.PostIssueCommentRequest{
		Owner:  "owner",
		Repo:   "repo",
		Number: 1,
		Token:  "fake-token",
		IssueComment: &scm.IssueComment{
			Body: &bodyStr,
		},
	}

	t.Run("Success_PostIssueComment", func(t *testing.T) {
		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "POST", req.Method)

				expectedPath := "/repos/owner/repo/issues/1/comments"
				assert.True(t, strings.Contains(req.URL.Path, expectedPath), "expected path to contain %s, got %s", expectedPath, req.URL.Path)

				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(strings.NewReader(`{"id": 9999}`)),
					Header:     make(http.Header),
				}, nil
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "token")
		require.NoError(t, err)

		err = driver.PostIssueComment(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("Failure_FailedToPostIssueComment", func(t *testing.T) {
		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`{"message": "Internal Server Error}`)),
					Header:     make(http.Header),
				}, nil
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "token")
		require.NoError(t, err)

		err = driver.PostIssueComment(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to post issue comment")
	})

	t.Run("Failure_NetworkError", func(t *testing.T) {
		transport := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("simulated network error")
			},
		}

		httpClient := &http.Client{Transport: transport}
		driver, err := scm.NewGitHubDriver(httpClient, "token")
		require.NoError(t, err)

		err = driver.PostIssueComment(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to post issue comment")
		assert.Contains(t, err.Error(), "simulated network error")
	})
}
