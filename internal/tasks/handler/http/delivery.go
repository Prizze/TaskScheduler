package http

import (
	"encoding/json"
	"net/http"

	"github.com/Prizze/TaskScheduler/internal/apperrors"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
	"github.com/Prizze/TaskScheduler/pkg/ctx"
	"github.com/Prizze/TaskScheduler/pkg/response"
)

type TasksHandler struct {
	service taskService
	cfg     *config.Config
}

func NewTasksHander(cfg *config.Config) *TasksHandler {
	return &TasksHandler{}
}

func (h *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ctx.UserIDKey).(int64)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	var req domain.CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.SendError(w, apperrors.Validation, nil)
		return
	}

	createTask := req.NewTask()
	taskWithTags, err := h.service.CreateTask(r.Context(), userID, createTask)
	if err != nil {
		handleError(w, err)
		return
	}

	resp := taskWithTags.NewTaskResponse()
	response.SendResponse(w, http.StatusCreated, resp)
}

func (h *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {

}

func (h *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {

}

func (h *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

}

func (h *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

}

func handleError(w http.ResponseWriter, err error) {

}
