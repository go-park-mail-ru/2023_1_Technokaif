package delivery

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) errorResponse(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	error := &errorResponse{Message: msg}
	message, err := json.Marshal(error)
	if err != nil {
		h.logger.Error("failed to marshal error message: " + err.Error())
	}

	w.Write([]byte(message))
}
