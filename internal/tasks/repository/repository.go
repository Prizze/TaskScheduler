package repository

import (
	"context"
	"database/sql"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/Prizze/TaskScheduler/internal/tasks/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TasksRepository struct {
	cfg  *config.Config
	pool *pgxpool.Pool
}

func NewTasksRepository(cfg *config.Config, pool *pgxpool.Pool) *TasksRepository {
	return &TasksRepository{
		cfg:  cfg,
		pool: pool,
	}
}

func (r *TasksRepository) CreateTask(ctx context.Context, userID int64, in *domain.CreateTask) (*domain.CreateTaskWithTags, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, domain.ErrServerError
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var res domain.CreateTaskWithTags

	tags := make([]*models.Tag, len(in.Tags))
	for i, t := range in.Tags {
		var tag models.Tag
		err := tx.QueryRow(ctx, getTagByID, t.ID, userID).Scan(&tag.ID, &tag.Name)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, domain.ErrNoTag
			}
			return nil, domain.ErrServerError
		}
		tags[i] = &tag
	}
	res.Tags = tags

	var dueDate sql.NullTime
	dueDateValue := any(nil)
	if !in.Task.DueDate.IsZero() {
		dueDateValue = in.Task.DueDate
	}

	task := &models.Task{}
	err = tx.QueryRow(
		ctx,
		createTask,
		userID,
		in.Task.Title,
		in.Task.Description,
		in.Task.Status,
		in.Task.Priority,
		dueDateValue,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&dueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, domain.ErrServerError
	}

	if dueDate.Valid {
		task.DueDate = dueDate.Time
	}

	for _, tag := range tags {
		if _, err := tx.Exec(ctx, createTaskTag, task.ID, tag.ID); err != nil {
			return nil, domain.ErrServerError
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, domain.ErrServerError
	}

	res.Task = task
	return &res, nil
}
