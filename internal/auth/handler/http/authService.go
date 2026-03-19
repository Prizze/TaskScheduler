package http

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/models"
)

//go:generate mockgen -source authService.go -destination=mocks/mock_service.go -package=mocks
type authService interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, user *models.User) (*models.User, error)
	Me(ctx context.Context, userID int64) (*models.User, error)
}
