package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const (
	LikeSuccess   = "ok"
	UnLikeSuccess = "ok"

	IncorrectRequestBody = "incorrect input body"
	InvalidURLParameter  = "invalid url parameter"
	UnathorizedUser      = "unathorized"

	SetLikeServerError    = "can't set like"
	DeleteLikeServerError = "can't remove like"

	LikeAlreadyExists = "already liked"
	LikeDoesntExist   = "wasn't liked"
)

type Error struct {
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, msg string, code int, logger logger.Logger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	errorResp := Error{Message: msg}
	message, err := json.Marshal(errorResp)
	if err != nil {
		logger.Errorf("failed to marshal error message: %v", err)
		if _, err = w.Write([]byte(`{"message": "can't encode error response into json"}`)); err != nil {
			logger.Errorf("failed to write response: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(message); err != nil {
		logger.Errorf("failed to write response: %v", err)
	}
}

const minErrorToLogCode = 500

func ErrorResponseWithErrLogging(w http.ResponseWriter, msg string, code int, logger logger.Logger, err error) {
	if err != nil {
		if code < minErrorToLogCode {
			logger.Info(err.Error())
		} else {
			logger.Error(err.Error())
		}
	}

	ErrorResponse(w, msg, code, logger)
}

func SuccessResponse(w http.ResponseWriter, r any, logger logger.Logger) {
	message, err := json.Marshal(r)
	if err != nil {
		logger.Error(err.Error())
		ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, logger)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(message); err != nil {
		logger.Errorf("failed to write response: %v", err)
	}
}
