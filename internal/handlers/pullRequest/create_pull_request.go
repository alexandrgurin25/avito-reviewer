package pull_request_handler

import (
	"avito-reviewer/internal/handlers"
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/models"
	"encoding/json"
	"net/http"
)

func (h *PRHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// не предусмотен такой ответ по API, но оставил для корректного ответа
		handlers.WriteBadRequest(w, r, "invalid json")
		return
	}

	in := models.PullRequest{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.AuthorID,
	}

	pr, err := h.s.CreatePullRequest(r.Context(), &in)
	if err != nil {
		handlers.WriteDomainError(w, r, err)
		return
	}

	out := dto.PullRequestDTO{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.Reviewers,
	}

	handlers.WriteJSON(w, r, http.StatusCreated, map[string]any{
		"pr": out,
	})
}
