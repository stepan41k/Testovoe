package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	musicHandler "github.com/stepan41k/Testovoe/internal/http-server/handlers/music"
	musicService "github.com/stepan41k/Testovoe/internal/service/music"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stepan41k/Testovoe/cmd/migrator"
	"github.com/stepan41k/Testovoe/internal/app"
	"github.com/stepan41k/Testovoe/internal/config"
	"github.com/stepan41k/Testovoe/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	storagePath := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username, cfg.Storage.DBName, os.Getenv("MY_DB_PASSWORD"), cfg.Storage.SSLMode)

	pool, err := postgres.New(context.Background(), storagePath)
	if err != nil {
		panic(err)
	}
	service := musicService.New(pool, log)
	handler := musicHandler.New(service, log)

	storagePathForMigrator := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", cfg.Storage.Username, os.Getenv("MY_DB_PASSWORD"), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.DBName, cfg.Storage.SSLMode)

	migrator.NewMigrator(storagePathForMigrator, os.Getenv("MY_MIGRATIONS_PATH"))

	router.Get("/songs", handler.GetSongs(context.Background()))

	router.Route("/song", func(r chi.Router) {
		r.Get("/text", handler.GetTextSong(context.Background()))
		r.Delete("/delete", handler.DeleteSong(context.Background()))
		r.Put("/update", handler.UpdateSong(context.Background()))
		r.Post("/new", handler.AddNewSong(context.Background()))
	})

	log.Info("starting server")

	application := app.New(log, cfg, router)

	go func() {
		application.HTTPServer.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <- stop

	log.Info("stopping application", slog.String("signal", signal.String()))

	application.HTTPServer.Stop(context.Background())

	postgres.Close(context.Background(), pool)

	log.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
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