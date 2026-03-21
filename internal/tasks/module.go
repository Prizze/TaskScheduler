package tasks

import (
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler int
	cfg     *config.Config
}

func NewTasksModule(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *Module {
	return &Module{
		
	}
}