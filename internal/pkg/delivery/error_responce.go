package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) errorResponce(w http.ResponseWriter, msg string, code int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	h.logger.Error(msg)

	error := &errorResponse{Message: msg}
	message, err := json.Marshal(error)
	if err != nil {
		return errors.New("failed to marshal error message")
	}

	w.Write([]byte(message))

	return nil
}
