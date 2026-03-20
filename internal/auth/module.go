package auth

import (
	"net/http"

	authhttp "github.com/Prizze/TaskScheduler/internal/auth/handler/http"
	"github.com/Prizze/TaskScheduler/internal/auth/repository"
	"github.com/Prizze/TaskScheduler/internal/auth/service"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/logger"
	"github.com/Prizze/TaskScheduler/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *authhttp.AuthHandler
	cfg     *config.Config
}

func NewAuthModule(db *pgxpool.Pool, cfg *config.Config, log logger.Logger) *Module {
	_ = log

	repo := repository.NewRepoAuth(db)
	service := service.NewAuthService(repo)
	handler := authhttp.NewAuthHandler(service, cfg)

	return &Module{
		cfg:     cfg,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", m.handler.Register)
	mux.HandleFunc("POST /auth/login", m.handler.Login)
	mux.Handle("GET /auth/me", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.Me)))
}
