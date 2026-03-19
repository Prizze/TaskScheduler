package middleware

import (
	"context"
	"net/http"

	"github.com/Prizze/TaskScheduler/internal/apperrors"
	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/Prizze/TaskScheduler/pkg/ctx"
	"github.com/Prizze/TaskScheduler/pkg/jwt"
	"github.com/Prizze/TaskScheduler/pkg/response"
)

func AuthHandler(cfg *config.Config,nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			response.SendError(w, apperrors.Unauthorized, nil)
			return 
		}

		claims, err := jwt.ParseJWT(cookie.Value, cfg)
		if err != nil {
			response.SendError(w, apperrors.Unauthorized, nil)
			return 
		}

		userID := claims.UserID

		ctx := context.WithValue(r.Context(), ctx.UserIDKey, userID)
		nextHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}