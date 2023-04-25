package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const realIPHeaderName = "X-Real-IP"

func Logging(logger logger.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			reqId := uuid.New()

			defer func() {
				respTime := time.Since(start)

				realIP := r.Header.Get(realIPHeaderName)
				if realIP != "" {
					logger.Infof("RID:%s, %s %s from RealIP: %s IP: %s - %s", reqId.String(), r.Method, r.URL.Path, realIP, r.RemoteAddr, respTime.String())
				} else {
					logger.Infof("RID:%s, %s %s from IP: %s - %s", reqId.String(), r.Method, r.URL.Path, r.RemoteAddr, respTime.String())
				}
			}()

			next.ServeHTTP(w, commonHttp.WrapReqID(r, reqId.ID()))
		})
	}
}
