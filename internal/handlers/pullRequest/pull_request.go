package pull_request_handler

import (
	"avito-reviewer/internal/services/pull_request_services"
)

type PRHandler struct {
	s pull_request_services.PRService
}

func NewPRHandler(s pull_request_services.PRService) *PRHandler {
	return &PRHandler{s: s}
}
