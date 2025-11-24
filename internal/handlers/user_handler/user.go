package userhand

import (
	userserv "avito-reviewer/internal/services/user_services"
)

type UserHandler struct {
	s userserv.UserService
}

func NewTeamHandler(s userserv.UserService) *UserHandler {
	return &UserHandler{s: s}
}
