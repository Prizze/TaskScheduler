package tags

import (
	"net/http"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/logger"
	"github.com/Prizze/TaskScheduler/internal/middleware"
	tagshttp "github.com/Prizze/TaskScheduler/internal/tags/handler/http"
	"github.com/Prizze/TaskScheduler/internal/tags/repository"
	"github.com/Prizze/TaskScheduler/internal/tags/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *tagshttp.TagsHandler
	cfg     *config.Config
}

func NewTagsModule(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *Module {
	repo := repository.NewTagsRepository(cfg, db)
	tagsService := service.NewTagsService(cfg, repo)
	handler := tagshttp.NewTagsHandler(tagsService, cfg)

	_ = log

	return &Module{
		handler: handler,
		cfg:     cfg,
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /tags", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.CreateTag)))
	mux.Handle("GET /tags", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.GetTags)))
	mux.Handle("DELETE /tags/{id}", middleware.AuthHandler(m.cfg, http.HandlerFunc(m.handler.DeleteTag)))
}
