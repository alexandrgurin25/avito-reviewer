package build

import (
	"avito-reviewer/internal/handlers"
	prhand "avito-reviewer/internal/handlers/pullRequest"
	teamhand "avito-reviewer/internal/handlers/team_handler"
	userhand "avito-reviewer/internal/handlers/user_handler"
	"net/http"

	"github.com/go-chi/chi"
)

func Router(team *teamhand.TeamHandler, pr *prhand.PRHandler,
	user *userhand.UserHandler) http.Handler {

	r := chi.NewRouter()

	r.Post("/team/add", team.AddTeam)
	r.Get("/team/get", team.GetTeam)

	r.Post("/users/setIsActive", user.SetIsActive)
	r.Get("/users/getReview", user.GetReview)

	r.Post("/pullRequest/create", pr.Create)
	r.Post("/pullRequest/merge", pr.Merge)
	r.Post("/pullRequest/reassign", pr.Reassign)

	// Дополнительный эндпоинт статистики -> количество назначений ревьюверов каждому пользователю
	r.Get("/stats/reviewers", pr.GetStats)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.WriteJSON(w, r, http.StatusOK, map[string]string{"status": "healthy"})
	})

	return r
}
