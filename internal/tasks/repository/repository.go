package repository

import (
	"context"
	"database/sql"
	"time"

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

type tagReader interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewTasksRepository(cfg *config.Config, pool *pgxpool.Pool) *TasksRepository {
	return &TasksRepository{
		cfg:  cfg,
		pool: pool,
	}
}

func (r *TasksRepository) CreateTask(ctx context.Context, userID int64, in *domain.CreateTask) (*domain.TaskWithTags, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, domain.ErrServerError
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	tags, err := r.validateTags(ctx, tx, userID, in.Tags)
	if err != nil {
		return nil, err
	}

	task, err := scanTaskRow(
		tx.QueryRow(
			ctx,
			createTask,
			userID,
			in.Task.Title,
			in.Task.Description,
			in.Task.Status,
			in.Task.Priority,
			nullableDueDate(in.Task.DueDate),
		),
	)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if _, err := tx.Exec(ctx, createTaskTag, task.ID, tag.ID); err != nil {
			return nil, domain.ErrServerError
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, domain.ErrServerError
	}

	return &domain.TaskWithTags{Task: task, Tags: tags}, nil
}

func (r *TasksRepository) GetTasks(ctx context.Context, userID int64) ([]*domain.TaskWithTags, error) {
	rows, err := r.pool.Query(ctx, getTasks, userID)
	if err != nil {
		return nil, domain.ErrServerError
	}
	defer rows.Close()

	var tasks []*domain.TaskWithTags
	for rows.Next() {
		task, scanErr := scanTaskRow(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		tags, tagsErr := r.fetchTaskTags(ctx, task.ID)
		if tagsErr != nil {
			return nil, tagsErr
		}

		tasks = append(tasks, &domain.TaskWithTags{Task: task, Tags: tags})
	}

	if rows.Err() != nil {
		return nil, domain.ErrServerError
	}

	return tasks, nil
}

func (r *TasksRepository) GetTask(ctx context.Context, userID, taskID int64) (*domain.TaskWithTags, error) {
	task, err := scanTaskRow(r.pool.QueryRow(ctx, getTaskByID, taskID, userID))
	if err != nil {
		return nil, err
	}

	tags, err := r.fetchTaskTags(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return &domain.TaskWithTags{Task: task, Tags: tags}, nil
}

func (r *TasksRepository) UpdateTask(ctx context.Context, userID, taskID int64, in *domain.UpdateTask) (*domain.TaskWithTags, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, domain.ErrServerError
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	tags, err := r.validateTags(ctx, tx, userID, in.Tags)
	if err != nil {
		return nil, err
	}

	task, err := scanTaskRow(
		tx.QueryRow(
			ctx,
			updateTask,
			taskID,
			userID,
			in.Task.Title,
			in.Task.Description,
			in.Task.Status,
			in.Task.Priority,
			nullableDueDate(in.Task.DueDate),
		),
	)
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, deleteTaskTags, taskID); err != nil {
		return nil, domain.ErrServerError
	}

	for _, tag := range tags {
		if _, err := tx.Exec(ctx, createTaskTag, taskID, tag.ID); err != nil {
			return nil, domain.ErrServerError
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, domain.ErrServerError
	}

	return &domain.TaskWithTags{Task: task, Tags: tags}, nil
}

func (r *TasksRepository) DeleteTask(ctx context.Context, userID, taskID int64) error {
	res, err := r.pool.Exec(ctx, deleteTask, taskID, userID)
	if err != nil {
		return domain.ErrServerError
	}
	if res.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func (r *TasksRepository) validateTags(ctx context.Context, db tagReader, userID int64, in []*models.Tag) ([]*models.Tag, error) {
	tags := make([]*models.Tag, len(in))
	for i, t := range in {
		var tag models.Tag
		err := db.QueryRow(ctx, getTagByID, t.ID, userID).Scan(&tag.ID, &tag.Name)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, domain.ErrNoTag
			}
			return nil, domain.ErrServerError
		}
		tags[i] = &tag
	}

	return tags, nil
}

func (r *TasksRepository) fetchTaskTags(ctx context.Context, taskID int64) ([]*models.Tag, error) {
	rows, err := r.pool.Query(ctx, getTaskTags, taskID)
	if err != nil {
		return nil, domain.ErrServerError
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, domain.ErrServerError
		}
		tags = append(tags, tag)
	}

	if rows.Err() != nil {
		return nil, domain.ErrServerError
	}

	return tags, nil
}

func scanTaskRow(row pgx.Row) (*models.Task, error) {
	var (
		task    models.Task
		dueDate sql.NullTime
	)

	err := row.Scan(
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
		if err == pgx.ErrNoRows {
			return nil, domain.ErrTaskNotFound
		}
		return nil, domain.ErrServerError
	}

	if dueDate.Valid {
		task.DueDate = dueDate.Time
	}

	return &task, nil
}

func nullableDueDate(dueDate time.Time) any {
	if !dueDate.IsZero() {
		return dueDate
	}

	return nil
}
