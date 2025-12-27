package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, we can't send a proper error response
		// since we've already written the status code
		_ = err
	}
}

func Error(w http.ResponseWriter, status int, message string, details interface{}) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    http.StatusText(status),
			Message: message,
			Details: details,
		},
	}
	JSON(w, status, response)
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, status, data)
}
