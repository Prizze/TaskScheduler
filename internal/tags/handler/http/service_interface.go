package http

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/models"
)

type tagService interface {
	CreateTag(ctx context.Context, userID int64, tag *models.Tag) (*models.Tag, error)
	GetTags(ctx context.Context, userID int64) ([]*models.Tag, error)
	DeleteTag(ctx context.Context, userID, tagID int64) error
}
