package userhand

import (
	"avito-reviewer/internal/handlers"
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/handlers/mappers"
	"net/http"
)

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("user_id")

	// В спецификации параметр user_id обязателен, поэтому добавил обработку, которая явна не прописана
	if userID == "" {
		handlers.WriteBadRequest(w, r, "user_id parameter is required")
		return
	}

	reviewers, err := h.s.GetReview(r.Context(), userID)

	if err != nil {
		handlers.WriteDomainError(w, r, err)
		return
	}

	handlers.WriteJSON(w, r, http.StatusOK, dto.GetReviewDTO{
		UserID:       reviewers.ID,
		PullRequests: mappers.MapReviewersToDTO(reviewers.PullRequests),
	})
}
