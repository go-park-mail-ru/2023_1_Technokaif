package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Error struct {
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, msg string, code int, logger logger.Logger) {
	errorResp := Error{Message: msg}
	message, err := json.Marshal(errorResp)
	if err != nil {
		logger.Errorf("failed to marshal error message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't encode error response into json, msg: " + msg))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(message)
}

func SuccessResponse(w http.ResponseWriter, r any, logger logger.Logger) {
	message, err := json.Marshal(r)
	if err != nil {
		logger.Error(err.Error())
		ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, logger)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(message)
}
