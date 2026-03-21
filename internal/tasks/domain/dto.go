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
	TagIDs      []int     `json:"tag_ids"`
}

func (req *CreateTaskRequest) NewTask() *models.Task {
	
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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
