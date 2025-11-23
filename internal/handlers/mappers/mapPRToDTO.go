package mappers

import (
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/models"
)

func PrToDTO(pr *models.PullRequest) dto.PullRequestDTO {

	return dto.PullRequestDTO{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.Reviewers,
	}
}
