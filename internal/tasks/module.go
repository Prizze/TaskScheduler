package tasks

import (
	"net/http"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/logger"
	taskshttp "github.com/Prizze/TaskScheduler/internal/tasks/handler/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler taskshttp.TasksHandler
	cfg     *config.Config
}

func NewTasksModule(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *Module {
	return &Module{}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /tasks", m.handler.CreateTask)
}
