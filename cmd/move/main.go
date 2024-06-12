package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hovanja2011/move/internal/config"
	"github.com/hovanja2011/move/internal/storage/postgresql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	log.Info("starting move", slog.String("env", cfg.Env))
	log.Debug("debug messanger are enabled")

	storage, err := postgresql.New(cfg.Storage.Path)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

// TODO: init router: chi

// TODO: run server:
