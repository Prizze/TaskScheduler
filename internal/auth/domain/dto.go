package domain

import (
	"time"

	"github.com/Prizze/TaskScheduler/internal/models"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (req *RegisterRequest) NewUser() *models.User {
	return &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}
}

type RegisterResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func RegisterResponseFromUser(user *models.User) *RegisterResponse {
	return &RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *LoginRequest) NewUser() *models.User {
	return &models.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

type LoginResponse struct {
	ID    int64  `json:"user_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func LoginResponseFromUser(user *models.User) *LoginResponse {
	return &LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}
