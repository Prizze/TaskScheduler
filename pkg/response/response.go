package response

import (
	"encoding/json"
	"net/http"

	"github.com/Prizze/TaskScheduler/internal/apperrors"
)

type errorResponse struct {
	Error errorInfo `json:"error"`
}

type errorInfo struct {
	Code    apperrors.Code `json:"code"`
	Message string         `json:"message"`
	Details any            `json:"details"`
}

func SendResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func SendError(w http.ResponseWriter, err apperrors.Response, details any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus)

	resp := errorResponse{
		Error: errorInfo{
			Code:    err.Code,
			Message: err.Message,
			Details: details,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
