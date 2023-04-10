package middleware

import (
	"net/http"
)

func RequestBodyMaxSize(maxSizeBytes uint64) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, int64(maxSizeBytes))
			next.ServeHTTP(w, r)
		})
	}
}
