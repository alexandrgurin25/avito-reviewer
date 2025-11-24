package teamhand

import "avito-reviewer/internal/services/team_services"

type TeamHandler struct {
	s teamserv.TeamService
}

func NewTeamHandler(s teamserv.TeamService) *TeamHandler {
	return &TeamHandler{s}
}
