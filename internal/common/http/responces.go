package http

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

	error := errorResponse{Message: msg}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(error); err != nil {
		logger.Errorf("failed to marshal error message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't encode error response into json, msg: " + msg))
	}
}

func SuccessResponse(w http.ResponseWriter, r any, logger logger.Logger) {
	w.Header().Set("Content-Type", "json/application; charset=utf-8")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(r); err != nil {
		logger.Error(err.Error())
		ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, logger)
	}
}
