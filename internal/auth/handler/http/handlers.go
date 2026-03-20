package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Prizze/TaskScheduler/internal/apperrors"
	"github.com/Prizze/TaskScheduler/internal/auth/domain"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/pkg/ctx"
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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.SendError(w, apperrors.Validation, nil)
		return
	}

	newUser := req.NewUser()

	user, err := h.Service.Register(r.Context(), newUser)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := setAuthCookie(w, user.ID, h.Cfg); err != nil {
		handleError(w, err)
		return
	}

	userDTO := domain.RegisterResponseFromUser(user)
	response.SendResponse(w, http.StatusCreated, userDTO)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
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

	if err := setAuthCookie(w, user.ID, h.Cfg); err != nil {
		handleError(w, err)
		return
	}

	userDTO := domain.LoginResponseFromUser(user)
	response.SendResponse(w, http.StatusOK, userDTO)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ctx.UserIDKey).(int64)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	user, err := h.Service.Me(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	userDTO := domain.RegisterResponseFromUser(user)
	response.SendResponse(w, http.StatusOK, userDTO)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		response.SendError(w, apperrors.Validation, err.Error())
		return
	case errors.Is(err, domain.ErrEmailIsTaken):
		response.SendError(w, apperrors.Conflict, domain.ErrEmailIsTaken.Error())
		return
	case errors.Is(err, domain.ErrInvalidCredentials):
		response.SendError(w, apperrors.Unauthorized, domain.ErrInvalidCredentials.Error())
		return
	case errors.Is(err, domain.ErrUserNotFound):
		response.SendError(w, apperrors.NotFound, domain.ErrUserNotFound.Error())
		return
	default:
		response.SendError(w, apperrors.Internal, nil)
		return
	}
}

func setAuthCookie(w http.ResponseWriter, userID int64, cfg *config.Config) error {
	token, err := jwt.GenerateJWT(userID, cfg)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}
