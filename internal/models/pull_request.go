package models

import "time"

type PullRequestStatus string

const (
	PROpen   PullRequestStatus = "OPEN"
	PRMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID        string
	Name      string
	AuthorID  string
	Status    PullRequestStatus
	CreatedAt time.Time
	MergedAt  *time.Time
	Reviewers []string
}

type UserPR struct {
	ID           string
	PullRequests []PullRequest
}
