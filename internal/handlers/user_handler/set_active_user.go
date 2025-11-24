package userhand

import (
	"avito-reviewer/internal/handlers"
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/models"
	"encoding/json"
	"net/http"
)

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req dto.TeamMemberDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Этот ответ не предусмотрен API,
		// но оставлен для обеспечения корректной работы.
		handlers.WriteBadRequest(w, r, "invalid json")
		return
	}

	team := models.User{
		ID:       req.UserID,
		IsActive: req.IsActive,
	}

	updated, err := h.s.SetIsActiveUser(r.Context(), &team)

	if err != nil {
		handlers.WriteDomainError(w, r, err)
		return
	}

	handlers.WriteJSON(w, r, http.StatusOK, map[string]any{
		"user": dto.TeamMemberDTO{
			UserID:   updated.ID,
			Username: updated.Username,
			TeamName: updated.TeamName,
			IsActive: updated.IsActive,
		},
	})
}
