package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hovanja2011/move/internal/config"
	mwLogger "github.com/hovanja2011/move/internal/http-server/middleware/logger"
	"github.com/hovanja2011/move/internal/lib/logger/handlers/slogpretty"
	sl "github.com/hovanja2011/move/internal/lib/logger/sl"
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

	log.Error("error messages are enabled")
	storage, err := postgresql.New(cfg.Storage.Path)
	if err != nil {
		log.Error("failed to init storage", sl.Error(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
