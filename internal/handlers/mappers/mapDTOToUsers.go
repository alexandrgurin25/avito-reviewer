package mappers

import (
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/models"
)

func MapDTOToUsers(members []dto.TeamMemberDTO) []models.User {
	res := make([]models.User, len(members))
	for i, m := range members {
		res[i] = models.User{
			ID:       m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return res
}
