package http

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
)

type taskService interface {
	CreateTask(ctx context.Context, task *models.Task, tags []*models.Tag) (*domain.CreateTaskWithTags, error)
}
