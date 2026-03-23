package app

import (
	"context"
	"net/http"
	"os"

	"github.com/Prizze/TaskScheduler/internal/auth"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/db"
	"github.com/Prizze/TaskScheduler/internal/logger"
	"github.com/Prizze/TaskScheduler/internal/middleware"
	"github.com/Prizze/TaskScheduler/internal/tags"
	"github.com/Prizze/TaskScheduler/internal/tasks"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	server *http.Server
	dbPool *pgxpool.Pool
}

func NewApp(cfg *config.Config) (*App, error) {
	mux := http.NewServeMux()

	logger := logger.NewSlogJSONLogger(os.Stdout, nil)
	dbPool, err := db.NewPostgresPool(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	authModule := auth.NewAuthModule(dbPool, cfg, logger)
	authModule.RegisterRoutes(mux)

	tasksModule := tasks.NewTasksModule(cfg, dbPool, logger)
	tasksModule.RegisterRoutes(mux)

	tagsModule := tags.NewTagsModule(cfg, dbPool, logger)
	tagsModule.RegisterRoutes(mux)

	return &App{
		server: &http.Server{
			Addr:    cfg.HTTPAddr,
			Handler: middleware.Recovery(logger.With("layer", "http"), mux),
		},
		dbPool: dbPool,
	}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Close() {
	if a.dbPool != nil {
		a.dbPool.Close()
	}
}
