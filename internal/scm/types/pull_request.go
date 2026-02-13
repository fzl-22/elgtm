package scm

import "time"

type PullRequest struct {
	ID        int64
	Number    int
	Title     string
	Body      string
	URL       string
	DiffURL   string
	RawDiff   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
