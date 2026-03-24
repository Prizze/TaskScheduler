package tasks

import (
	"net/http"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/logger"
	"github.com/Prizze/TaskScheduler/internal/middleware"
	taskshttp "github.com/Prizze/TaskScheduler/internal/tasks/handler/http"
	"github.com/Prizze/TaskScheduler/internal/tasks/repository"
	"github.com/Prizze/TaskScheduler/internal/tasks/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *taskshttp.TasksHandler
	cfg     *config.Config
}

func NewTasksModule(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *Module {
	repo := repository.NewTasksRepository(cfg, db)
	tasksService := service.NewTasksService(cfg, repo)
	handler := taskshttp.NewTasksHander(tasksService, cfg)

	_ = log

	return &Module{
		handler: handler,
		cfg:     cfg,
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /tasks", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.CreateTask)))
	mux.Handle("GET /tasks", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.GetTasks)))
	mux.Handle("GET /tasks/{id}", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.GetTask)))
	mux.Handle("PUT /tasks/{id}", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.UpdateTask)))
	mux.Handle("DELETE /tasks/{id}", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.DeleteTask)))
}
