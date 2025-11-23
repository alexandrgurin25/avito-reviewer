package handlers

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/pkg/logger"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type apiErrorResponse struct {
	Error apiErrorDetail `json:"error"`
}

type apiErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeAPIError(w http.ResponseWriter, r *http.Request, code, message string, status int) {
	ctx := r.Context()

	log := logger.GetLoggerFromCtx(ctx)

	log.Info(ctx,
		"API error",

		zap.Any("status", status),
		zap.Any("method", r.Method),
		zap.Any("path", r.URL.Path),
		zap.Any("code", code),
		zap.Any("message", message),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := apiErrorResponse{
		Error: apiErrorDetail{
			Code:    code,
			Message: message,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func WriteBadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	writeAPIError(w, r, "BAD_REQUEST", msg, http.StatusBadRequest)
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, body any) {

	ctx := r.Context()

	log := logger.GetLoggerFromCtx(ctx)

	log.Info(ctx,
		"API response",
		zap.Any("status", status),
		zap.Any("method", r.Method),
		zap.Any("path", r.URL.Path),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(body)
}

func WriteDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case models.ErrPRExists:
		writeAPIError(w, r, "PR_EXISTS", err.Error(), http.StatusConflict)
	case models.ErrPRMerged:
		writeAPIError(w, r, "PR_MERGED", err.Error(), http.StatusConflict)
	case models.ErrNotAssigned:
		writeAPIError(w, r, "NOT_ASSIGNED", err.Error(), http.StatusConflict)
	case models.ErrNoCandidate:
		writeAPIError(w, r, "NO_CANDIDATE", err.Error(), http.StatusConflict)
	case models.ErrNotFound:
		writeAPIError(w, r, "NOT_FOUND", err.Error(), http.StatusNotFound)
	case models.ErrUserBelongsToAnotherTeam:
		writeAPIError(w, r, "USER_EXISTS", err.Error(), http.StatusBadRequest)
	case models.ErrTeamExists:
		writeAPIError(w, r, "TEAM_EXISTS", err.Error(), http.StatusBadRequest)
	default:
		writeAPIError(w, r, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
	}
}
