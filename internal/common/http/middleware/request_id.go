package middleware

import (
	"net/http"

	"github.com/google/uuid"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
)

func SetReqId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New()
		next.ServeHTTP(w, commonHttp.WrapReqID(r, reqId.ID()))
	})
}
