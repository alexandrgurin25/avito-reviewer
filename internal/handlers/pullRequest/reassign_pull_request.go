package pull_request_handler

import (
	"avito-reviewer/internal/handlers"
	"avito-reviewer/internal/handlers/dto"
	"avito-reviewer/internal/handlers/mappers"
	"avito-reviewer/internal/models"
	"encoding/json"
	"net/http"
)

func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req dto.ReassignPRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// не предусмотен такой ответ по API, но оставил для корректного ответа
		handlers.WriteBadRequest(w, r, "invalid json")
		return
	}

	in := models.ReasignPR{
		PRID:          req.PRID,
		OldReviewerID: req.OldReviewer,
	}

	pr, replacedBy, err := h.s.ReassignReviewer(r.Context(), &in)
	if err != nil {
		handlers.WriteDomainError(w, r, err)
		return
	}

	handlers.WriteJSON(w, r, http.StatusOK, dto.ReassignPRResponse{
		PR:         mappers.PrToDTO(pr),
		ReplacedBy: replacedBy,
	})
}
