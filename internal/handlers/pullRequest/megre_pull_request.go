package pull_request_handler

import (
	"avito-reviewer/internal/handlers"

	"avito-reviewer/internal/handlers/dto"
	"encoding/json"
	"net/http"
)

func (h *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	var req dto.MergePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// не предусмотен такой ответ по API, поэтому убрал
		// handlers.WriteBadRequest(w, r, "invalid json")
		return
	}

	pr, err := h.s.MergePR(r.Context(), req.ID)
	if err != nil {
		handlers.WriteDomainError(w, r, err)
		return
	}

	timeMerge := pr.MergedAt.UTC().Format("2006-01-02T15:04:05.000Z")

	out := dto.PullRequestDTO{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.Reviewers,
		MergedAt:          &timeMerge,
	}

	handlers.WriteJSON(w, r, http.StatusOK, map[string]any{
		"pr": out,
	})
}
