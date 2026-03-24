package service

//go:generate mockgen -source repository_interface.go -destination=mocks/mock_repo.go -package=mocks

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
)

type tasksRepository interface {
	CreateTask(ctx context.Context, userID int64, in *domain.CreateTask) (*domain.TaskWithTags, error)
	GetTasks(ctx context.Context, userID int64) ([]*domain.TaskWithTags, error)
	GetTask(ctx context.Context, userID, taskID int64) (*domain.TaskWithTags, error)
	UpdateTask(ctx context.Context, userID, taskID int64, in *domain.UpdateTask) (*domain.TaskWithTags, error)
	DeleteTask(ctx context.Context, userID, taskID int64) error
}
