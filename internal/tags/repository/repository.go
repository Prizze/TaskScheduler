package repository

import (
	"context"
	"errors"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/Prizze/TaskScheduler/internal/tags/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagsRepository struct {
	cfg  *config.Config
	pool *pgxpool.Pool
}

func NewTagsRepository(cfg *config.Config, pool *pgxpool.Pool) *TagsRepository {
	return &TagsRepository{
		cfg:  cfg,
		pool: pool,
	}
}

func (r *TagsRepository) CreateTag(ctx context.Context, userID int64, tag *models.Tag) (*models.Tag, error) {
	createdTag := &models.Tag{}
	err := r.pool.QueryRow(ctx, createTag, userID, tag.Name).Scan(
		&createdTag.ID,
		&createdTag.UserID,
		&createdTag.Name,
		&createdTag.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrTagAlreadyExists
		}
		return nil, domain.ErrServerError
	}

	return createdTag, nil
}

func (r *TagsRepository) GetTags(ctx context.Context, userID int64) ([]*models.Tag, error) {
	rows, err := r.pool.Query(ctx, getTags, userID)
	if err != nil {
		return nil, domain.ErrServerError
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, domain.ErrServerError
		}
		tags = append(tags, tag)
	}

	if rows.Err() != nil {
		return nil, domain.ErrServerError
	}

	return tags, nil
}

func (r *TagsRepository) DeleteTag(ctx context.Context, userID, tagID int64) error {
	res, err := r.pool.Exec(ctx, deleteTag, tagID, userID)
	if err != nil {
		return domain.ErrServerError
	}
	if res.RowsAffected() == 0 {
		return domain.ErrTagNotFound
	}

	return nil
}
