package dto

type TeamRequest struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name,omitempty"`
	IsActive bool   `json:"is_active"`
}

type TeamResponse struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type CreatePRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type MergePRRequest struct {
	ID string `json:"pull_request_id"`
}

type ReassignPRRequest struct {
	PRID        string `json:"pull_request_id"`
	OldReviewer string `json:"old_user_id"`
}

type ReassignPRResponse struct {
	PR         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}

type PullRequestDTO struct {
	ID                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers,omitempty"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type GetReviewDTO struct {
	UserID       string           `json:"user_id"`
	PullRequests []PullRequestDTO `json:"pull_requests"`
}
