package response

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error errorInfo `json:"error"`
}

type errorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func SendResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func SendError(w http.ResponseWriter, status int, code string, message string, details any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	err := errorResponse{
		Error: errorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}

	_ = json.NewEncoder(w).Encode(err)
}
