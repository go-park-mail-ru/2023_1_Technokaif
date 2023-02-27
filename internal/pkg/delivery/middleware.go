package delivery

import (
	"context"
	"net/http"
	"strings"
)

// HTTP middleware setting a value on the request context
func (h Handler) Authorization(next http.Handler) http.Handler { // TEST IT
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer "
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)

		if authHeader == "" || reqToken == authHeader {
			httpErrorResponce(w, "no token found", http.StatusBadRequest)
			return
		}
		userId, err := h.services.CheckAccessToken(reqToken)
		if err != nil {
			httpErrorResponce(w, "failed to authenticate", http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), "ID", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
