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

func NewTasksService(cfg *config.Config, repo tasksRepository) *TasksService {
	return &TasksService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *TasksService) CreateTask(ctx context.Context, userID int64, input *domain.CreateTask) (*domain.TaskWithTags, error) {
	if err := validateTask(input.Task); err != nil {
		return nil, err
	}

	createdTask, err := s.repo.CreateTask(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	return enrichTask(createdTask), nil
}

func (s *TasksService) GetTasks(ctx context.Context, userID int64) ([]*domain.TaskWithTags, error) {
	tasks, err := s.repo.GetTasks(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		enrichTask(task)
	}

	return tasks, nil
}

func (s *TasksService) GetTask(ctx context.Context, userID, taskID int64) (*domain.TaskWithTags, error) {
	task, err := s.repo.GetTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	return enrichTask(task), nil
}

func (s *TasksService) UpdateTask(ctx context.Context, userID, taskID int64, input *domain.UpdateTask) (*domain.TaskWithTags, error) {
	if err := validateTask(input.Task); err != nil {
		return nil, err
	}

	updatedTask, err := s.repo.UpdateTask(ctx, userID, taskID, input)
	if err != nil {
		return nil, err
	}

	return enrichTask(updatedTask), nil
}

func (s *TasksService) DeleteTask(ctx context.Context, userID, taskID int64) error {
	return s.repo.DeleteTask(ctx, userID, taskID)
}

func enrichTask(taskWithTags *domain.TaskWithTags) *domain.TaskWithTags {
	if taskWithTags != nil && taskWithTags.Task != nil {
		taskWithTags.IsOverdue = isTaskOverdue(taskWithTags.Task)
	}

	return taskWithTags
}

func validateTask(task *models.Task) error {
	if task == nil || task.Title == "" {
		return domain.ErrValidation
	}

	switch task.Status {
	case models.StatusPending, models.StatusInProgress, models.StatusDone:
	default:
		return domain.ErrInvalidStatus
	}

	switch task.Priority {
	case models.PriorityLow, models.PriorityMedium, models.PriorityHigh:
	default:
		return domain.ErrInvalidPriority
	}

	return nil
}

func isTaskOverdue(task *models.Task) bool {
	if task == nil || task.DueDate.IsZero() {
		return false
	}
	if task.Status == models.StatusDone {
		return false
	}

	return true
}
