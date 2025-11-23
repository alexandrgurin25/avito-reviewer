package team_services

import (
	"avito-reviewer/internal/models"
	"context"
)

func (s teamService) AddTeam(ctx context.Context, t *models.Team) (*models.Team, error) {

	//Начинаем транзакцию
	tx, err := s.teamRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	//Проверка, что такой команды еще нет
	hasTeam, err := s.teamRepo.TeamExists(ctx, tx, t.Name)

	if err != nil {
		return nil, err
	}

	if hasTeam {
		return nil, models.ErrTeamExists
	}

	// Так как  API не предусматривает вариант того, что один пользователь может быть только в одной команде и не может менять ее, то эту реализацию скрываем
	// p.s. с моей стороны логично, что условно один разработчик работает только в одной команде, но буду следовать требованиям :)

	// // Проверка, что таких участников еще нет, а если есть, то их команда соответствует
	// userIDs := make([]string, len(t.Members))
	// for _, k := range t.Members {
	// 	userIDs = append(userIDs, k.ID)
	// }

	// existing, err := s.userRepo.GetExistingUsers(ctx, tx, userIDs)

	// for _, m := range t.Members {
	// 	if team, ok := existing[m.ID]; ok {
	// 		if team != t.Name {
	// 			return nil, models.ErrUserBelongsToAnotherTeam
	// 		}
	// 	}
	// }

	// if err != nil {
	// 	return nil, err
	// }

	// Создаем команду
	createdTeam, err := s.teamRepo.CreateTeam(ctx, tx, t.Name)

	if err != nil {
		return nil, err
	}

	t.ID = createdTeam.ID

	//Вносим батчами юзеров
	createdUsers, err := s.userRepo.UpsertUsers(ctx, tx, t)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	createdTeam.Members = createdUsers.Members

	return createdTeam, nil

}
