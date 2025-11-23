package main

import (
	"avito-reviewer/internal/config"
	pull_request_handler "avito-reviewer/internal/handlers/pullRequest"
	teamhand "avito-reviewer/internal/handlers/team_handler"
	userhand "avito-reviewer/internal/handlers/user_handler"
	"avito-reviewer/internal/repositories"
	prRepo "avito-reviewer/internal/repositories/pull_request_repository"
	teamRepo "avito-reviewer/internal/repositories/team_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"

	prserv "avito-reviewer/internal/services/pull_request_services"
	teamserv "avito-reviewer/internal/services/team_services"
	"avito-reviewer/internal/services/user_services"
	"avito-reviewer/pkg/logger"
	"avito-reviewer/pkg/postgres"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	log := logger.GetLoggerFromCtx(ctx)

	cfg, err := config.New()
	if err != nil {
		log.Fatal(ctx, "unable to load config", zap.Error(err))
		return
	}

	DB, err := postgres.New(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "unable to connect db", zap.Error(err))
		return
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Successful start!")

	r := chi.NewRouter()
	repository := repositories.NewPgxPoolAdapter(DB)

	userRepository := userrepo.NewUserRepository(repository)
	teamRepository := teamRepo.NewTeamRepository(repository)
	prRepository := prRepo.NewPRRepository(repository)

	teamServices := teamserv.NewService(userRepository, teamRepository, repository)
	userServices := userserv.NewService(userRepository, teamRepository, prRepository, repository)
	prServices := prserv.NewService(userRepository, prRepository)

	teamHandler := teamhand.NewTeamHandler(teamServices)
	userHandler := userhand.NewTeamHandler(userServices)
	prHandler := pull_request_handler.NewPRHandler(prServices)

	r.Route("/team/", func(r chi.Router) {
		r.Post("/add", teamHandler.AddTeam)
		r.Get("/get", teamHandler.GetTeam)
	})

	r.Route("/users/", func(r chi.Router) {
		r.Post("/setIsActive", userHandler.SetIsActive)
		r.Get("/getReview", userHandler.GetReview)
	})

	r.Route("/pullRequest/", func(r chi.Router) {
		r.Post("/create", prHandler.Create)
		r.Post("/reassign", prHandler.Reassign)
		r.Post("/merge", prHandler.Merge)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Error(ctx, "Server error", zap.Error(err))
		return

	case <-shutdown:
		log.Info(ctx, "Starting graceful shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error(ctx, "Graceful shutdown failed", zap.Error(err))
			server.Close()
		}
	}

	log.Info(ctx, "Server stopped")
}
