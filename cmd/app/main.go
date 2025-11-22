package main

import (
	"avito-reviewer/internal/config"
	"avito-reviewer/internal/handlers/team_handler"
	"avito-reviewer/internal/handlers/user_handler"
	"avito-reviewer/internal/repositories"
	"avito-reviewer/internal/repositories/team_repository"
	"avito-reviewer/internal/repositories/user_repository"
	"avito-reviewer/internal/services/team_services"
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

	userRepository := user_repository.NewUserRepository(repository)
	teamRepository := team_repository.NewTeamRepository(repository)

	teamServices := team_services.NewService(userRepository, teamRepository, repository)
	userServices := user_services.NewService(userRepository, teamRepository, repository)

	teamHandler := team_handler.NewTeamHandler(teamServices)
	userHandler := user_handler.NewTeamHandler(userServices)

	r.Route("/team/", func(r chi.Router) {
		r.Post("/add/", teamHandler.AddTeam)
		r.Get("/get/", teamHandler.GetTeam)
	})

	r.Route("/users/", func(r chi.Router) {
		r.Post("/setIsActive/", userHandler.SetIsActive)
		// r.Get("/get/", teamHandler.GetTeam)
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
