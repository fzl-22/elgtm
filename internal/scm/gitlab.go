package scm

import (
	"fmt"

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
