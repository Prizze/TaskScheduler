package http

//go:generate mockgen -source service_interface.go -destination=mocks/mock_service.go -package=mocks

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
)

type taskService interface {
	CreateTask(ctx context.Context, userID int64, in *domain.CreateTask) (*domain.TaskWithTags, error)
	GetTask(ctx context.Context, userID, taskID int64) (*domain.TaskWithTags, error)
	UpdateTask(ctx context.Context, userID, taskID int64, in *domain.UpdateTask) (*domain.TaskWithTags, error)
	DeleteTask(ctx context.Context, userID, taskID int64) error
}
