package delivery

import (
	"net/http"
	"context"
)

// HTTP middleware setting a value on the request context
func (h *Handler) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  ctx := context.WithValue(r.Context(), "user", "123")
	  next.ServeHTTP(w, r.WithContext(ctx))
	})
  }