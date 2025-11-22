package mappers

import (
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/models"
)

func MapUsersToDTO(users []models.User) []dto.TeamMemberDTO {
	res := make([]dto.TeamMemberDTO, len(users))
	for i, u := range users {
		res[i] = dto.TeamMemberDTO{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		}
	}
	return res
}
