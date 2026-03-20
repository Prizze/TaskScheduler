package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Prizze/TaskScheduler/internal/auth/domain"
	servicemocks "github.com/Prizze/TaskScheduler/internal/auth/service/mocks"
	loggermocks "github.com/Prizze/TaskScheduler/internal/logger/mocks"
	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthService_PanicsOnNilLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := servicemocks.NewMockrepoAuth(ctrl)

	assert.Panics(t, func() {
		NewAuthService(repo, nil)
	})
}

func TestAuthService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMockrepoAuth(ctrl)
		log := loggermocks.NewMockLogger(ctrl)
		service := NewAuthService(repo, log)

		input := &models.User{
			Email:    "  user@example.com  ",
			Password: "password123",
			Name:     "  Ivan  ",
		}

		repo.EXPECT().CheckEmailExists(gomock.Any(), "user@example.com").Return(false, nil)
		repo.EXPECT().
			CreateUser(gomock.Any(), gomock.AssignableToTypeOf(&models.User{})).
			DoAndReturn(func(_ context.Context, user *models.User) (*models.User, error) {
				require.Equal(t, "user@example.com", user.Email)
				require.Equal(t, "Ivan", user.Name)
				require.NotEmpty(t, user.PasswordHash)
				require.NotEqual(t, input.Password, user.PasswordHash)

				return &models.User{
					ID:           42,
					Email:        user.Email,
					Name:         user.Name,
					PasswordHash: user.PasswordHash,
				}, nil
			})
		log.EXPECT().Info("auth register succeeded", "user_id", int64(42), "email", "user@example.com")

		user, err := service.Register(context.Background(), input)

		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, int64(42), user.ID)
		assert.Equal(t, "user@example.com", user.Email)
		assert.Equal(t, "Ivan", user.Name)
	})

	t.Run("email taken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMockrepoAuth(ctrl)
		log := loggermocks.NewMockLogger(ctrl)
		service := NewAuthService(repo, log)

		input := &models.User{
			Email:    "user@example.com",
			Password: "password123",
			Name:     "Ivan",
		}

		repo.EXPECT().CheckEmailExists(gomock.Any(), "user@example.com").Return(true, nil)
		log.EXPECT().Info("auth register email already taken", "email", "user@example.com")

		user, err := service.Register(context.Background(), input)

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrEmailIsTaken)
		assert.Nil(t, user)
	})

	t.Run("validation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMockrepoAuth(ctrl)
		log := loggermocks.NewMockLogger(ctrl)
		service := NewAuthService(repo, log)

		log.EXPECT().Warn("auth register validation failed")

		user, err := service.Register(context.Background(), &models.User{
			Email:    "",
			Password: "short",
			Name:     "",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrValidation)
		assert.Nil(t, user)
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMockrepoAuth(ctrl)
		log := loggermocks.NewMockLogger(ctrl)
		service := NewAuthService(repo, log)

		hashedPassword, err := hashPassword("password123")
		require.NoError(t, err)

		repo.EXPECT().GetUserByEmail(gomock.Any(), "user@example.com").Return(&models.User{
			ID:           7,
			Email:        "user@example.com",
			Name:         "Ivan",
			PasswordHash: hashedPassword,
		}, nil)
		log.EXPECT().Info("auth login succeeded", "user_id", int64(7), "email", "user@example.com")

		user, err := service.Login(context.Background(), &models.User{
			Email:    " user@example.com ",
			Password: "password123",
		})

		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, int64(7), user.ID)
		assert.Equal(t, "user@example.com", user.Email)
	})

	t.Run("invalid credentials when user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMockrepoAuth(ctrl)
		log := loggermocks.NewMockLogger(ctrl)
		service := NewAuthService(repo, log)

		repo.EXPECT().GetUserByEmail(gomock.Any(), "user@example.com").Return(nil, nil)
		log.EXPECT().Info("auth login user not found", "email", "user@example.com")

		user, err := service.Login(context.Background(), &models.User{
			Email:    "user@example.com",
			Password: "password123",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Nil(t, user)
	})

	t.Run("repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMockrepoAuth(ctrl)
		log := loggermocks.NewMockLogger(ctrl)
		service := NewAuthService(repo, log)

		repo.EXPECT().GetUserByEmail(gomock.Any(), "user@example.com").Return(nil, errors.New("db error"))
		log.EXPECT().Error("auth login get user failed", "email", "user@example.com", "err", gomock.Any())

		user, err := service.Login(context.Background(), &models.User{
			Email:    "user@example.com",
			Password: "password123",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrServerError)
		assert.Nil(t, user)
	})
}

func TestAuthService_Me(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMockrepoAuth(ctrl)
		log := loggermocks.NewMockLogger(ctrl)
		service := NewAuthService(repo, log)

		expectedUser := &models.User{
			ID:    1,
			Email: "user@example.com",
			Name:  "Ivan",
		}

		repo.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(expectedUser, nil)

		user, err := service.Me(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := servicemocks.NewMockrepoAuth(ctrl)
		log := loggermocks.NewMockLogger(ctrl)
		service := NewAuthService(repo, log)

		repo.EXPECT().GetUserByID(gomock.Any(), int64(10)).Return(nil, nil)
		log.EXPECT().Info("auth me user not found", "user_id", int64(10))

		user, err := service.Me(context.Background(), 10)

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, user)
	})
}
