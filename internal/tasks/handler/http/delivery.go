package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

func NewTasksHander(service taskService, cfg *config.Config) *TasksHandler {
	return &TasksHandler{
		service: service,
		cfg:     cfg,
	}
}

func (h *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	var req domain.CreateTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.SendError(w, apperrors.Validation, nil)
		return
	}

	taskWithTags, err := h.service.CreateTask(r.Context(), userID, req.NewTask())
	if err != nil {
		handleError(w, err)
		return
	}

	response.SendResponse(w, http.StatusCreated, taskWithTags.NewTaskResponse())
}

func (h *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	tasks, err := h.service.GetTasks(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.SendResponse(w, http.StatusOK, domain.TasksResponseFromModels(tasks))
}

func (h *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	taskID, err := taskIDFromRequest(r)
	if err != nil {
		response.SendError(w, apperrors.Validation, err.Error())
		return
	}

	taskWithTags, err := h.service.GetTask(r.Context(), userID, taskID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.SendResponse(w, http.StatusOK, taskWithTags.NewTaskResponse())
}

func (h *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	taskID, err := taskIDFromRequest(r)
	if err != nil {
		response.SendError(w, apperrors.Validation, err.Error())
		return
	}

	var req domain.UpdateTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.SendError(w, apperrors.Validation, nil)
		return
	}

	taskWithTags, err := h.service.UpdateTask(r.Context(), userID, taskID, req.NewTask())
	if err != nil {
		handleError(w, err)
		return
	}

	response.SendResponse(w, http.StatusOK, taskWithTags.NewTaskResponse())
}

func (h *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	taskID, err := taskIDFromRequest(r)
	if err != nil {
		response.SendError(w, apperrors.Validation, err.Error())
		return
	}

	if err := h.service.DeleteTask(r.Context(), userID, taskID); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation),
		errors.Is(err, domain.ErrInvalidStatus),
		errors.Is(err, domain.ErrInvalidPriority):
		response.SendError(w, apperrors.Validation, err.Error())
	case errors.Is(err, domain.ErrNoTag),
		errors.Is(err, domain.ErrTaskNotFound):
		response.SendError(w, apperrors.NotFound, err.Error())
	default:
		response.SendError(w, apperrors.Internal, nil)
	}
}

func userIDFromContext(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(ctx.UserIDKey).(int64)
	return userID, ok
}

func taskIDFromRequest(r *http.Request) (int64, error) {
	taskID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || taskID <= 0 {
		return 0, domain.ErrValidation
	}

	return taskID, nil
}
