package app

import (
	"log/slog"

	"github.com/go-chi/chi"
	httpapp "github.com/stepan41k/Testovoe/internal/app/http"
	"github.com/stepan41k/Testovoe/internal/config"
)

type App struct {
	HTTPServer *httpapp.App
	log *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config, router chi.Router) *App {
	
	httpApp :=	httpapp.New(log, cfg, router)
	
	return &App{
		HTTPServer: httpApp,
		log: log,
	}
}

