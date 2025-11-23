package userserv

import (
	"avito-reviewer/internal/models"
	"context"
)

func (s *userService) GetReview(ctx context.Context, userID string) (*models.UserPR, error) {

	//Находим все pr_id
	prIDs, err := s.prRepo.GetIDByReviewerID(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}
	//Достаем по pr_id все проверки

	pullRequests, err := s.prRepo.GetPRsByID(ctx, s.db, prIDs)

	if err != nil {
		return nil, err
	}

	res := models.UserPR{
		ID:           userID,
		PullRequests: pullRequests,
	}

	return &res, nil
}
