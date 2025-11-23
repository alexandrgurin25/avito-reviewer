package team_handler

import "avito-reviewer/internal/services/team_services"

type TeamHandler struct {
	s team_services.TeamService
}

func NewTeamHandler(s team_services.TeamService) *TeamHandler {
	return &TeamHandler{s}
}
