package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	easyjson "github.com/mailru/easyjson"
)

//go:generate easyjson -no_std_marshalers responses.go

const (
	LikeSuccess   = "ok"
	UnLikeSuccess = "ok"

	IncorrectRequestBody = "incorrect input body"
	InvalidURLParameter  = "invalid url parameter"
	UnathorizedUser      = "unathorized"
	ForbiddenUser        = "user has no rights"

	SetLikeServerError    = "can't set like"
	DeleteLikeServerError = "can't remove like"

	LikeAlreadyExists = "already liked"
	LikeDoesntExist   = "wasn't liked"
)

//easyjson:json
type Error struct {
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, msg string, code int, logger logger.Logger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	errorResp := Error{Message: msg}
	message, err := easyjson.Marshal(errorResp)
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

func ErrorResponseWithErrLogging(w http.ResponseWriter, r *http.Request,
	msg string, code int, logger logger.Logger, err error) {

	if err != nil {
		if code < minErrorToLogCode {
			logger.InfoReqID(r.Context(), err.Error())
		} else {
			logger.ErrorReqID(r.Context(), err.Error())
		}
	}

	ErrorResponse(w, r, msg, code, logger)
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, response easyjson.Marshaler, logger logger.Logger) {
	message, err := easyjson.Marshal(response)
	if err != nil {
		logger.Error(err.Error())
		ErrorResponseWithErrLogging(w, r, "can't encode response into json", http.StatusInternalServerError, logger, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(message); err != nil {
		logger.Errorf("failed to write response: %v", err)
	}
}
