package service

import (
	"context"
	"testing"
	"time"

	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
	servicemocks "github.com/Prizze/TaskScheduler/internal/tasks/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTasksService_CreateTask(t *testing.T) {
	t.Run("success sets overdue flag", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMocktasksRepository(ctrl)
		service := NewTasksService(nil, repo)

		input := &domain.CreateTask{
			Task: &models.Task{
				Title:    "Pay bills",
				Status:   models.StatusPending,
				Priority: models.PriorityHigh,
				DueDate:  time.Now().Add(-time.Hour),
			},
		}

		repo.EXPECT().CreateTask(gomock.Any(), int64(7), input).Return(&domain.TaskWithTags{
			Task: &models.Task{
				ID:       10,
				Title:    input.Task.Title,
				Status:   input.Task.Status,
				Priority: input.Task.Priority,
				DueDate:  input.Task.DueDate,
			},
		}, nil)

		result, err := service.CreateTask(context.Background(), 7, input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsOverdue)
	})

	t.Run("validation error when title is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMocktasksRepository(ctrl)
		service := NewTasksService(nil, repo)

		result, err := service.CreateTask(context.Background(), 7, &domain.CreateTask{
			Task: &models.Task{
				Title:    "",
				Status:   models.StatusPending,
				Priority: models.PriorityHigh,
			},
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrValidation)
		assert.Nil(t, result)
	})
}

func TestTasksService_GetTask(t *testing.T) {
	t.Run("done task is not overdue", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMocktasksRepository(ctrl)
		service := NewTasksService(nil, repo)

		repo.EXPECT().GetTask(gomock.Any(), int64(5), int64(9)).Return(&domain.TaskWithTags{
			Task: &models.Task{
				ID:       9,
				Title:    "Done task",
				Status:   models.StatusDone,
				Priority: models.PriorityLow,
				DueDate:  time.Now().Add(-2 * time.Hour),
			},
		}, nil)

		result, err := service.GetTask(context.Background(), 5, 9)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsOverdue)
	})
}

func TestTasksService_UpdateTask(t *testing.T) {
	t.Run("invalid priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMocktasksRepository(ctrl)
		service := NewTasksService(nil, repo)

		result, err := service.UpdateTask(context.Background(), 1, 2, &domain.UpdateTask{
			Task: &models.Task{
				Title:    "Task",
				Status:   models.StatusPending,
				Priority: models.Priority("urgent"),
			},
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidPriority)
		assert.Nil(t, result)
	})
}

func TestTasksService_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := servicemocks.NewMocktasksRepository(ctrl)
	service := NewTasksService(nil, repo)

	repo.EXPECT().DeleteTask(gomock.Any(), int64(3), int64(11)).Return(nil)

	err := service.DeleteTask(context.Background(), 3, 11)

	require.NoError(t, err)
}
