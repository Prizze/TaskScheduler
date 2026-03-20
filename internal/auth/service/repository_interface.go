package service

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/models"
)

//go:generate mockgen -source repository_interface.go -destination=mocks/mock_repo.go -package=mocks
type repoAuth interface {
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}
