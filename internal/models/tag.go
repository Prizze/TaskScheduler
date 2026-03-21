package models

import "time"

type Tag struct {
	ID        int64
	UserID    int64
	Name      string
	CreatedAt time.Time
}
