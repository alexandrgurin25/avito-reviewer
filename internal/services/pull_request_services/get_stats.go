package prserv

import "context"

func (s *prService) GetReviewerStats(ctx context.Context) (map[string]int, error) {
	return s.prRepo.GetReviewerStats(ctx, s.db)
}
