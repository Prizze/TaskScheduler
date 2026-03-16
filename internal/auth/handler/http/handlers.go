package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Prizze/TaskScheduler/internal/apperrors"
	"github.com/Prizze/TaskScheduler/internal/auth/domain"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/pkg/jwt"
	"github.com/Prizze/TaskScheduler/pkg/response"
)

type AuthHandler struct {
	Service authService
	Cfg     *config.Config
}

func NewAuthHandler(service authService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		Service: service,
		Cfg:     cfg,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, apperrors.Validation, nil)
		return
	}

	newUser := req.NewUser()

	user, err := h.Service.Register(r.Context(), newUser)
	if err != nil {
		handleError(w, err)
		return
	}

	token, err := jwt.GenerateJWT(user.ID, h.Cfg)
	if err != nil {
		handleError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	userDTO := domain.RegisterResponseFromUser(user)
	response.SendResponse(w, http.StatusCreated, userDTO)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.SendError(w, apperrors.Validation, nil)
		return
	}

	userLogin := req.NewUser()
	user, err := h.Service.Login(r.Context(), userLogin)
	if err != nil {
		handleError(w, err)
		return
	}

	token, err := jwt.GenerateJWT(user.ID, h.Cfg)
	if err != nil {
		handleError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	userDTO := domain.LoginResponseFromUser(user)
	response.SendResponse(w, http.StatusOK, userDTO)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailIsTaken):
		response.SendError(w, apperrors.Conflict, domain.ErrEmailIsTaken.Error())
		return
	}
}
