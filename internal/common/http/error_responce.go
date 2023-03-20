package common_http

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	error := &errorResponse{Message: msg}
	message, _ := json.Marshal(error)

	// TODO
	/* if err != nil {
		h.logger.Errorf("failed to marshal error message: %s", err.Error())
	} */

	w.Write([]byte(message))
}
