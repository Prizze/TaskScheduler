package service

import (
	"context"
	"strings"

	"github.com/Prizze/TaskScheduler/internal/auth/domain"
	"github.com/Prizze/TaskScheduler/internal/models"
)

type AuthService struct {
	repo repoAuth
}

func NewAuthService(repo repoAuth) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateRegisterUser(user); err != nil {
		return nil, err
	}

	exist, err := s.isEmailTaken(ctx, user.Email)
	if err != nil {
		return nil, domain.ErrServerError
	}
	if exist {
		return nil, domain.ErrEmailIsTaken
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrServerError
	}

	userToCreate := &models.User{
		Email:        strings.TrimSpace(user.Email),
		PasswordHash: hashedPassword,
		Name:         strings.TrimSpace(user.Name),
	}

	createdUser, err := s.repo.CreateUser(ctx, userToCreate)
	if err != nil {
		return nil, domain.ErrServerError
	}

	return createdUser, nil
}

func (s *AuthService) Login(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateLoginUser(user); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.GetUserByEmail(ctx, strings.TrimSpace(user.Email))
	if err != nil {
		return nil, domain.ErrServerError
	}
	if existingUser == nil {
		return nil, domain.ErrInvalidCredentials
	}

	matched, err := verifyPassword(user.Password, existingUser.PasswordHash)
	if err != nil {
		return nil, domain.ErrServerError
	}
	if !matched {
		return nil, domain.ErrInvalidCredentials
	}

	return existingUser, nil
}

func (s *AuthService) Me(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrServerError
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (s *AuthService) isEmailTaken(ctx context.Context, email string) (bool, error) {
	return s.repo.CheckEmailExists(ctx, strings.TrimSpace(email))
}

func validateRegisterUser(user *models.User) error {
	if user == nil {
		return domain.ErrValidation
	}
	if strings.TrimSpace(user.Email) == "" {
		return domain.ErrValidation
	}
	if strings.TrimSpace(user.Name) == "" {
		return domain.ErrValidation
	}
	if len(user.Password) < 8 {
		return domain.ErrValidation
	}

	return nil
}

func validateLoginUser(user *models.User) error {
	if user == nil {
		return domain.ErrValidation
	}
	if strings.TrimSpace(user.Email) == "" {
		return domain.ErrValidation
	}
	if user.Password == "" {
		return domain.ErrValidation
	}

	return nil
}
