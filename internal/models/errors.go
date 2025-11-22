package models

import "errors"

var (
	ErrTeamExists  = errors.New("team_name already exists")
	ErrPRExists    = errors.New("pull request already exists")
	ErrPRMerged    = errors.New("pull request is already merged")
	ErrNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
	ErrNotFound    = errors.New("resource not found")
	ErrUserBelongsToAnotherTeam = errors.New("user belongs to another team")
)
