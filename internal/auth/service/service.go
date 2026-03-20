package service

import (
	"context"
	"strings"

	"github.com/Prizze/TaskScheduler/internal/auth/domain"
	"github.com/Prizze/TaskScheduler/internal/logger"
	"github.com/Prizze/TaskScheduler/internal/models"
)

type AuthService struct {
	repo   repoAuth
	logger logger.Logger
}

func NewAuthService(repo repoAuth, log logger.Logger) *AuthService {
	if log == nil {
		panic("auth service logger is nil")
	}

	return &AuthService{
		repo:   repo,
		logger: log,
	}
}

func (s *AuthService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateRegisterUser(user); err != nil {
		s.logger.Warn("auth register validation failed")
		return nil, err
	}

	exist, err := s.isEmailTaken(ctx, user.Email)
	if err != nil {
		s.logger.Error("auth register check email failed", "email", strings.TrimSpace(user.Email), "err", err)
		return nil, domain.ErrServerError
	}
	if exist {
		s.logger.Info("auth register email already taken", "email", strings.TrimSpace(user.Email))
		return nil, domain.ErrEmailIsTaken
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		s.logger.Error("auth register hash password failed", "email", strings.TrimSpace(user.Email), "err", err)
		return nil, domain.ErrServerError
	}

	userToCreate := &models.User{
		Email:        strings.TrimSpace(user.Email),
		PasswordHash: hashedPassword,
		Name:         strings.TrimSpace(user.Name),
	}

	createdUser, err := s.repo.CreateUser(ctx, userToCreate)
	if err != nil {
		s.logger.Error("auth register create user failed", "email", userToCreate.Email, "err", err)
		return nil, domain.ErrServerError
	}

	s.logger.Info("auth register succeeded", "user_id", createdUser.ID, "email", createdUser.Email)

	return createdUser, nil
}

func (s *AuthService) Login(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateLoginUser(user); err != nil {
		s.logger.Warn("auth login validation failed")
		return nil, err
	}

	email := strings.TrimSpace(user.Email)

	existingUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("auth login get user failed", "email", email, "err", err)
		return nil, domain.ErrServerError
	}
	if existingUser == nil {
		s.logger.Info("auth login user not found", "email", email)
		return nil, domain.ErrInvalidCredentials
	}

	matched, err := verifyPassword(user.Password, existingUser.PasswordHash)
	if err != nil {
		s.logger.Error("auth login verify password failed", "email", email, "err", err)
		return nil, domain.ErrServerError
	}
	if !matched {
		s.logger.Info("auth login invalid password", "email", email)
		return nil, domain.ErrInvalidCredentials
	}

	s.logger.Info("auth login succeeded", "user_id", existingUser.ID, "email", existingUser.Email)

	return existingUser, nil
}

func (s *AuthService) Me(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("auth me get user failed", "user_id", userID, "err", err)
		return nil, domain.ErrServerError
	}
	if user == nil {
		s.logger.Info("auth me user not found", "user_id", userID)
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
