package middleware

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const realIPHeaderName = "X-Real-IP"

func Logging(logger logger.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				respTime := time.Since(start)
				realIP := r.Header.Get(realIPHeaderName)
				if realIP != "" {
					logger.InfofReqID(r.Context(), "%s %s from RealIP: %s IP: %s - %s", r.Method, r.URL.Path, realIP, r.RemoteAddr, respTime.String())
				} else {
					logger.InfofReqID(r.Context(), "%s %s from IP: %s - %s", r.Method, r.URL.Path, r.RemoteAddr, respTime.String())
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
