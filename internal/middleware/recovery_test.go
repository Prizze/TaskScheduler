package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	loggermocks "github.com/Prizze/TaskScheduler/internal/logger/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	t.Run("passes request through", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := loggermocks.NewMockLogger(ctrl)
		handler := Recovery(log, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("recovers panic and returns internal error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := loggermocks.NewMockLogger(ctrl)
		log.EXPECT().Error(
			"http panic recovered",
			"method", http.MethodGet,
			"path", "/panic",
			"panic", "boom",
			"stack", gomock.Any(),
		)

		handler := Recovery(log, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("boom")
		}))

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), `"code":"internal_error"`)
	})
}
