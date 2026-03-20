package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Prizze/TaskScheduler/internal/apperrors"
	"github.com/Prizze/TaskScheduler/internal/logger"
	"github.com/Prizze/TaskScheduler/pkg/response"
)

func Recovery(log logger.Logger, next http.Handler) http.Handler {
	if log == nil {
		panic("recovery middleware logger is nil")
	}
	if next == nil {
		panic("recovery middleware handler is nil")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error(
					"http panic recovered",
					"method", r.Method,
					"path", r.URL.Path,
					"panic", fmt.Sprint(recovered),
					"stack", string(debug.Stack()),
				)

				response.SendError(w, apperrors.Internal, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
