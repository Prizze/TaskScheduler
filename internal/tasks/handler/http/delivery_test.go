package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
	"github.com/Prizze/TaskScheduler/internal/tasks/handler/http/mocks"
	pkgctx "github.com/Prizze/TaskScheduler/pkg/ctx"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTasksHandler_GetTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mocks.NewMocktaskService(ctrl)
		handler := NewTasksHander(service, &config.Config{})
		req := requestWithUser(http.MethodGet, "/tasks/15", nil, 42)
		req.SetPathValue("id", "15")
		rec := httptest.NewRecorder()

		service.EXPECT().GetTask(req.Context(), int64(42), int64(15)).Return(sampleTaskWithTags(), nil)

		handler.GetTask(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp domain.TaskResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		assert.Equal(t, int64(15), resp.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mocks.NewMocktaskService(ctrl)
		handler := NewTasksHander(service, &config.Config{})
		req := requestWithUser(http.MethodGet, "/tasks/bad", nil, 42)
		req.SetPathValue("id", "bad")
		rec := httptest.NewRecorder()

		handler.GetTask(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mocks.NewMocktaskService(ctrl)
		handler := NewTasksHander(service, &config.Config{})
		req := requestWithUser(http.MethodGet, "/tasks/99", nil, 42)
		req.SetPathValue("id", "99")
		rec := httptest.NewRecorder()

		service.EXPECT().GetTask(req.Context(), int64(42), int64(99)).Return(nil, domain.ErrTaskNotFound)

		handler.GetTask(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestTasksHandler_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMocktaskService(ctrl)
	handler := NewTasksHander(service, &config.Config{})
	body := domain.UpdateTaskRequest{
		Title:       "Updated title",
		Description: "Updated description",
		Status:      string(models.StatusInProgress),
		Priority:    string(models.PriorityMedium),
		DueDate:     time.Now().UTC().Truncate(time.Second),
		TagIDs:      []int64{1, 2},
	}
	req := requestWithUser(http.MethodPut, "/tasks/15", body, 42)
	req.SetPathValue("id", "15")
	rec := httptest.NewRecorder()

	service.EXPECT().UpdateTask(req.Context(), int64(42), int64(15), body.NewTask()).Return(sampleTaskWithTags(), nil)

	handler.UpdateTask(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTasksHandler_DeleteTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mocks.NewMocktaskService(ctrl)
		handler := NewTasksHander(service, &config.Config{})
		req := requestWithUser(http.MethodDelete, "/tasks/15", nil, 42)
		req.SetPathValue("id", "15")
		rec := httptest.NewRecorder()

		service.EXPECT().DeleteTask(req.Context(), int64(42), int64(15)).Return(nil)

		handler.DeleteTask(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mocks.NewMocktaskService(ctrl)
		handler := NewTasksHander(service, &config.Config{})
		req := httptest.NewRequest(http.MethodDelete, "/tasks/15", nil)
		req.SetPathValue("id", "15")
		rec := httptest.NewRecorder()

		handler.DeleteTask(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func requestWithUser(method, target string, data any, userID int64) *http.Request {
	var body *bytes.Buffer
	if data != nil {
		payload, _ := json.Marshal(data)
		body = bytes.NewBuffer(payload)
	} else {
		body = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, target, body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), pkgctx.UserIDKey, userID)
	return req.WithContext(ctx)
}

func sampleTaskWithTags() *domain.TaskWithTags {
	now := time.Now().UTC().Truncate(time.Second)
	return &domain.TaskWithTags{
		Task: &models.Task{
			ID:          15,
			Title:       "Task title",
			Description: "Task description",
			Status:      models.StatusPending,
			Priority:    models.PriorityHigh,
			DueDate:     now.Add(time.Hour),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Tags: []*models.Tag{
			{ID: 1, Name: "backend"},
		},
	}
}
