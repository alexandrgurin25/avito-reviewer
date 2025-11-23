package prhand

import (
	"avito-reviewer/internal/services/pull_request_services"
)

type PRHandler struct {
	s prserv.PRService
}

func NewPRHandler(s prserv.PRService) *PRHandler {
	return &PRHandler{s: s}
}
