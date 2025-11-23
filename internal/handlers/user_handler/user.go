package user_handler

import (
	"avito-reviewer/internal/services/user_services"
)

type UserHandler struct {
	s user_services.UserService
}

func NewTeamHandler(s user_services.UserService) *UserHandler {
	return &UserHandler{s: s}
}
