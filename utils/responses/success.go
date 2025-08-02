package responses

import (
	"encoding/json"
	"net/http"
)

type Success struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteSuccessResponse(w http.ResponseWriter, message string, data any) {
	writeSuccess(w, message, data, http.StatusOK)
}

func WriteSuccessResponseWithCode(w http.ResponseWriter, message string, data any, code int) {
	writeSuccess(w, message, data, code)
}

func writeSuccess(w http.ResponseWriter, message string, data any, code int) {
	resp := Success{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
