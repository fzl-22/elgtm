package scm

import "time"

type PullRequest struct {
	ID        int64
	Number    int
	Title     string
	Body      string
	Author    string
	URL       string
	DiffURL   string
	RawDiff   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IssueComment struct {
	Body *string
}
