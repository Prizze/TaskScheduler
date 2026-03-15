package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SomeData struct {
	Num    int    `json:"num"`
	String string `json:"string"`
}

func TestSendResponse(t *testing.T) {
	t.Run("Status OK", func(t *testing.T) {
		resp := httptest.NewRecorder()
		data := SomeData{
			Num:    1,
			String: "111",
		}

		SendResponse(resp, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, resp.Header().Get("Content-Type"), "application/json")

		var respData SomeData
		_ = json.NewDecoder(resp.Body).Decode(&respData)

		assert.Equal(t, data, respData)
	})
}

func TestSendError(t *testing.T) {
	resp := httptest.NewRecorder()

	err := ErrorResponse{
		HTTPStatus: http.StatusBadRequest,
		Code:       "bad request",
		Message:    "invalid request data",
	}

	SendError(resp, err, nil)

	response := errorResponse{}

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, resp.Header().Get("Content-Type"), "application/json")

	_ = json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, err.Code, response.Error.Code)
	assert.Equal(t, err.Message, response.Error.Message)
}
