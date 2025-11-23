package user_services

import (
	"avito-reviewer/internal/models"
	"context"
)

func (s *userService) GetReview(ctx context.Context, userID string) (*models.UserPR, error) {

	//Находим все pr_id
	prIDs, err := s.prRepo.GetIDByReviewerId(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}
	//Достаем по pr_id все проверки

	pullRequests, err := s.prRepo.GetPRsById(ctx, s.db, prIDs)

	if err != nil {
		return nil, err
	}

	res := models.UserPR{
		ID:           userID,
		PullRequests: pullRequests,
	}

	return &res, nil
}
