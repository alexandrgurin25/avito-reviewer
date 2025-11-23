package userserv

import (
	"avito-reviewer/internal/models"
	"context"
)

func (s *userService) SetIsActiveUser(ctx context.Context, u *models.User) (*models.User, error) {

	hasUser, err := s.userRepo.UserExists(ctx, s.db, u.ID)
	if err != nil {
		return nil, err
	}

	if !hasUser {
		return nil, models.ErrNotFound
	}

	resp, teamID, err := s.userRepo.SetActive(ctx, s.db, u.ID, u.IsActive)
	if err != nil {
		return nil, err
	}

	team, err := s.teamRepo.GetTeamNameByID(ctx, s.db, teamID)
	if err != nil {
		return nil, err
	}

	resp.TeamName = team.Name

	return resp, nil

}
