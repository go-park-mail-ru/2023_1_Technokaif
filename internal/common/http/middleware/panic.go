package middleware

import (
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

func Panic(logger logger.Logger) func(next http.Handler) http.Handler {
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("PANIC: %s", err)
					logger.Info("panic recovered")
					commonHttp.ErrorResponse(w, "server panic", http.StatusInternalServerError, logger)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}