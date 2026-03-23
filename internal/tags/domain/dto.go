package domain

import (
	"strings"
	"time"

	"github.com/Prizze/TaskScheduler/internal/models"
)

type CreateTagRequest struct {
	Name string `json:"name"`
}

func (req *CreateTagRequest) NewTag() *models.Tag {
	return &models.Tag{
		Name: strings.TrimSpace(req.Name),
	}
}

type TagResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func TagResponseFromModel(tag *models.Tag) *TagResponse {
	if tag == nil {
		return nil
	}

	return &TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
	}
}

type TagsResponse struct {
	Tags []TagResponse `json:"tags"`
}

func TagsResponseFromModels(tags []*models.Tag) *TagsResponse {
	items := make([]TagResponse, len(tags))
	for i, tag := range tags {
		items[i] = TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
		}
	}

	return &TagsResponse{Tags: items}
}
