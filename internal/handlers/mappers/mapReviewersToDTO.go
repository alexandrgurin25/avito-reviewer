package mappers

import (
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/models"
)

func MapReviewersToDTO(prs []models.PullRequest) []dto.PullRequestDTO {
	res := make([]dto.PullRequestDTO, len(prs))
	for i, p := range prs {
		res[i] = dto.PullRequestDTO{
			ID:       p.ID,
			Name:     p.Name,
			AuthorID: p.AuthorID,
			Status:   string(p.Status),
		}
	}
	return res
}
