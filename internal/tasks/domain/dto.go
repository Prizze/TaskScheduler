package domain

import (
	"time"

	"github.com/Prizze/TaskScheduler/internal/models"
)

type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	TagIDs      []int64   `json:"tag_ids"`
}

func (req *CreateTaskRequest) NewTask() (*models.Task, []*models.Tag) {
	tags := make([]*models.Tag, len(req.TagIDs))
	for i, id := range req.TagIDs {
		tag := &models.Tag{
			ID: id,
		}
		tags[i] = tag
	}

	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}, tags
}

type CreateTaskWithTags struct {
	Task      *models.Task
	Tags      []*models.Tag
	IsOverdue bool
}

func (res *CreateTaskWithTags) NewTaskResponse() *TaskResponse {
	tags := make([]Tag, len(res.Tags))
	for i, tag := range res.Tags {
		tags[i] = Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return &TaskResponse{
		ID:          res.Task.ID,
		Title:       res.Task.Title,
		Description: res.Task.Description,
		Status:      res.Task.Status,
		Priority:    res.Task.Priority,
		DueDate:     res.Task.DueDate,
		CreatedAt:   res.Task.CreatedAt,
		UpdatedAt:   res.Task.CreatedAt,
		IsOverdue:   res.IsOverdue,
		Tags:        tags,
	}
}

type TaskResponse struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      models.Status   `json:"status"`
	Priority    models.Priority `json:"priority"`
	DueDate     time.Time       `json:"due_date"`
	Tags        []Tag           `json:"tags"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"upadted_at"`
	IsOverdue   bool            `json:"is_overdue"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
