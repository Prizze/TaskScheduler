package models

import "time"

// Priority - Приоритеты задач
type Priority string

const PriorityHigh Priority = "high"
const PriorityMedium Priority = "medium"
const PriorityLow Priority = "low"

// Status - Статус задачи
type Status string

const StatusPending = "pending"
const StatusInProgress = "in_progress"
const StatusDone = "done"

type Task struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	StatusID    int
	Status      Status
	PriorityID  int
	Priority    Priority
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
