package service

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
)

type tasksRepository interface {
	CreateTask(ctx context.Context, userID int64, in *domain.CreateTask) (*domain.CreateTaskWithTags, error)
}
