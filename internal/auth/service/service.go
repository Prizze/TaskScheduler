package service

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/models"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (s *AuthService) Login(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (s *AuthService) Me(ctx context.Context, userID int64) (*models.User, error) {
	return nil, nil
}
