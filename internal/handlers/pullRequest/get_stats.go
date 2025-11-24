package prhand

import (
	"avito-reviewer/internal/handlers"
	"net/http"
)

func (h *PRHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.s.GetReviewerStats(r.Context())
	if err != nil {
		handlers.WriteDomainError(w, r, err)
		return
	}

	handlers.WriteJSON(w, r, 200, map[string]any{
		"reviewer_stats": stats,
	})
}
