package app

import (
	"net/http"
	"os"

	"github.com/Prizze/TaskScheduler/internal/auth"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/logger"
)

type App struct {
	server *http.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	mux := http.NewServeMux()

	_ = cfg

	logger := logger.NewSlogJSONLogger(os.Stdout, nil)

	authModule := auth.NewAuthModule(nil, cfg, logger)
	authModule.RegisterRoutes(mux)

	return &App{
		server: &http.Server{
			Addr:    cfg.HTTPAddr,
			Handler: mux,
		},
	}, nil
}
