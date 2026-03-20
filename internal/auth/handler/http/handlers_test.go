package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Prizze/TaskScheduler/internal/auth/domain"
	"github.com/Prizze/TaskScheduler/internal/auth/handler/http/mocks"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func getRegisterRequest() domain.RegisterRequest {
	return domain.RegisterRequest{
		Email:    "some@email.ru",
		Name:     "Ivan",
		Password: "password123",
	}
}

func getUser() *models.User {
	return &models.User{
		ID:    1,
		Email: "some@email.ru",
		Name:  "Ivan",
	}
}

func getRequestWithBody(method string, target string, data any) *http.Request {
	body, _ := json.Marshal(data)
	req := httptest.NewRequest(method, target, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockauthService(ctrl)
	cfg := &config.Config{JwtSecret: "secret"}
	handler := NewAuthHandler(mockService, cfg)

	t.Run("Register ok", func(t *testing.T) {
		dto := getRegisterRequest()
		user := getUser()
		req := getRequestWithBody(http.MethodPost, "/register", dto)
		rec := httptest.NewRecorder()

		mockService.EXPECT().Register(req.Context(), dto.NewUser()).Return(user, nil)

		handler.Register(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockauthService(ctrl)
	cfg := &config.Config{JwtSecret: "secret"}
	handler := NewAuthHandler(mockService, cfg)

	t.Run("Login ok", func(t *testing.T) {
		dto := domain.LoginRequest{
			Email:    "some@email.ru",
			Password: "password123",
		}

		req := getRequestWithBody(http.MethodPost, "/login", dto)
		rec := httptest.NewRecorder()

		mockService.EXPECT().Login(req.Context(), dto.NewUser()).Return(getUser(), nil)

		handler.Login(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
