package service

import (
	"context"
	"strings"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/Prizze/TaskScheduler/internal/tags/domain"
)

type TagsService struct {
	cfg  *config.Config
	repo tagsRepository
}

func NewTagsService(cfg *config.Config, repo tagsRepository) *TagsService {
	return &TagsService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *TagsService) CreateTag(ctx context.Context, userID int64, tag *models.Tag) (*models.Tag, error) {
	if err := validateTag(tag); err != nil {
		return nil, err
	}

	return s.repo.CreateTag(ctx, userID, tag)
}

func (s *TagsService) GetTags(ctx context.Context, userID int64) ([]*models.Tag, error) {
	return s.repo.GetTags(ctx, userID)
}

func (s *TagsService) DeleteTag(ctx context.Context, userID, tagID int64) error {
	return s.repo.DeleteTag(ctx, userID, tagID)
}

func validateTag(tag *models.Tag) error {
	if tag == nil {
		return domain.ErrValidation
	}

	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		return domain.ErrValidation
	}

	return nil
}
