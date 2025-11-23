package team_handler

import (
	"avito-reviewer/internal/handlers"
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/handlers/mappers"
	"avito-reviewer/internal/models"
	"encoding/json"
	"net/http"
)

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.TeamResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// не предусмотен такой ответ по API, поэтому убрал
		// handlers.WriteBadRequest(w, r, "invalid json")
		return
	}

	if req.TeamName == "" || len(req.Members) == 0 {
		// не предусмотен такой ответ по API, поэтому убрал
		// handlers.WriteBadRequest(w, r, "invalid json")
		return
	}

	team := models.Team{
		Name:    req.TeamName,
		Members: mappers.MapDTOToUsers(req.Members),
	}

	created, err := h.s.AddTeam(r.Context(), &team)

	if err != nil {
		handlers.WriteDomainError(w, r, err)
		return
	}

	handlers.WriteJSON(w, r, http.StatusCreated, dto.TeamResponse{
		TeamName: created.Name,
		Members:  mappers.MapUsersToDTO(created.Members),
	})
}
