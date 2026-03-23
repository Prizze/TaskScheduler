package service

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
)

type TasksService struct {
	cfg  *config.Config
	repo tasksRepository
}

func NewTasksService(cfg *config.Config) *TasksService {
	return &TasksService{
		cfg: cfg,
	}
}

func (s *TasksService) CreateTask(ctx context.Context, userID int64, input *domain.CreateTask) (*domain.CreateTaskWithTags, error) {
	createdTask, err := s.repo.CreateTask(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	if createdTask != nil && createdTask.Task != nil {
		createdTask.IsOverdue = isTaskOverdue(createdTask.Task)
	}

	return createdTask, nil
}

func isTaskOverdue(task *models.Task) bool {
	if task == nil || task.DueDate.IsZero() {
		return false
	}
	if task.Status == models.Status("done") {
		return false
	}
}
