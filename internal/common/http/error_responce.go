package common_http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type errorResponse struct {
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, msg string, code int, logger logger.Logger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	error := &errorResponse{Message: msg}
	message, err := json.Marshal(error)

	if err != nil {
		logger.Errorf("failed to marshal error message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(message))
}
