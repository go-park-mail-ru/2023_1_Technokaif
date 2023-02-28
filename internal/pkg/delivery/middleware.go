package delivery

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// HTTP middleware setting a value on the request context
func (h *Handler) Authorization(next http.Handler) http.Handler { // TEST IT
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer "
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)

		h.logger.Info("auth token : " + reqToken)

		if authHeader == "" || reqToken == authHeader {
			next.ServeHTTP(w, r) // missing token
			return
		}

		userId, err := h.services.CheckAccessToken(reqToken)
		if err != nil {
			h.logger.Error(err.Error())
			next.ServeHTTP(w, r) // token check failed
			return
		}

		ctx := context.WithValue(r.Context(), "ID", userId)
		next.ServeHTTP(w, r.WithContext(ctx)) // token check successed
	})
}

// returns error if authentication failed
func (h *Handler) GetIdFromAuthorization(r *http.Request) (uint, error) {
	userId, ok := r.Context().Value("ID").(uint)
	if !ok {
		h.logger.Error("no User ID")
		return 0, errors.New("failed to authenticate")
	}

	return userId, nil
}
