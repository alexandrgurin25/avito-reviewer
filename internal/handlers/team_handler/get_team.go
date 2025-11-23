package team_handler

import (
	"avito-reviewer/internal/handlers"
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/handlers/mappers"
	"net/http"
)

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	teamName := r.URL.Query().Get("team_name")

	// По API не требовалось обрабатывать  
	// if teamName == "" {
	// 	handlers.WriteBadRequest(w, r, "empty query parameters")
	// 	return
	// }

	foundTeam, err := h.s.GetTeam(ctx, teamName)

	if err != nil {
		handlers.WriteDomainError(w, r, err)
		return
	}

	handlers.WriteJSON(w, r, http.StatusOK, dto.TeamResponse{
		TeamName: foundTeam.Name,
		Members:  mappers.MapUsersToDTO(foundTeam.Members),
	})
}
