package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Prizze/TaskScheduler/internal/apperrors"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/internal/tags/domain"
	pkgctx "github.com/Prizze/TaskScheduler/pkg/ctx"
	"github.com/Prizze/TaskScheduler/pkg/response"
)

type TagsHandler struct {
	service tagService
	cfg     *config.Config
}

func NewTagsHandler(service tagService, cfg *config.Config) *TagsHandler {
	return &TagsHandler{
		service: service,
		cfg:     cfg,
	}
}

func (h *TagsHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	var req domain.CreateTagRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.SendError(w, apperrors.Validation, nil)
		return
	}

	tag, err := h.service.CreateTag(r.Context(), userID, req.NewTag())
	if err != nil {
		handleError(w, err)
		return
	}

	response.SendResponse(w, http.StatusCreated, domain.TagResponseFromModel(tag))
}

func (h *TagsHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	tags, err := h.service.GetTags(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.SendResponse(w, http.StatusOK, domain.TagsResponseFromModels(tags))
}

func (h *TagsHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.SendError(w, apperrors.Unauthorized, nil)
		return
	}

	tagID, err := tagIDFromRequest(r)
	if err != nil {
		response.SendError(w, apperrors.Validation, err.Error())
		return
	}

	if err := h.service.DeleteTag(r.Context(), userID, tagID); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		response.SendError(w, apperrors.Validation, err.Error())
	case errors.Is(err, domain.ErrTagAlreadyExists):
		response.SendError(w, apperrors.Conflict, err.Error())
	case errors.Is(err, domain.ErrTagNotFound):
		response.SendError(w, apperrors.NotFound, err.Error())
	default:
		response.SendError(w, apperrors.Internal, nil)
	}
}

func userIDFromContext(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(pkgctx.UserIDKey).(int64)
	return userID, ok
}

func tagIDFromRequest(r *http.Request) (int64, error) {
	tagID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || tagID <= 0 {
		return 0, domain.ErrValidation
	}

	return tagID, nil
}
