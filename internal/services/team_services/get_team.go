package teamserv

import (
	"avito-reviewer/internal/models"
	"context"
)

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {

	// Проверяем есть ли такая команда
	hasTeam, err := s.teamRepo.TeamExists(ctx, s.db, teamName)
	if err != nil {
		return nil, err
	}

	if !hasTeam {
		return nil, models.ErrNotFound
	}

	team, err := s.teamRepo.GetTeamIDByName(ctx, s.db, teamName)
	if err != nil {
		return nil, err
	}

	response, err := s.userRepo.GetUsersByTeam(ctx, s.db, team.ID)
	if err != nil {
		return nil, err
	}

	response.Name = teamName

	return response, nil
}
