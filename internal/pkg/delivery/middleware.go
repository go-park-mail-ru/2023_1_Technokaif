package delivery

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

const contextValueUser = "user"

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

		userId, userVersion, err := h.services.CheckAccessToken(reqToken)
		if err != nil {
			h.logger.Error(err.Error())
			next.ServeHTTP(w, r) // token check failed
			return
		}

		user, err := h.services.GetUserByAuthData(userId, userVersion)
		if err != nil {
			h.logger.Error(err.Error())
			next.ServeHTTP(w, r) // UserAuth data check failed
			return
		}

		ctx := context.WithValue(r.Context(), contextValueUser, user)
		next.ServeHTTP(w, r.WithContext(ctx)) // token check successed
	})
}

// returns error if authentication failed
func (h *Handler) GetUserFromAuthorization(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(contextValueUser).(*models.User)
	if !ok {
		h.logger.Error("no User")
		return nil, errors.New("failed to authenticate")
	}

	return user, nil
}
