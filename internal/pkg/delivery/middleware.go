package delivery

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type contextKeyUserType struct{}

// Authorization is HTTP middleware which sets a value on the request context
func (h *Handler) authorization(next http.Handler) http.Handler { // TEST IT
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer"
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)
		reqToken = strings.ReplaceAll(reqToken, " ", "")

		h.logger.Info("auth token : " + reqToken)

		if authHeader == "" || reqToken == authHeader || reqToken == "" {
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

		h.logger.Info("user version : " + strconv.Itoa(int(user.Version)))

		ctx := context.WithValue(r.Context(), contextKeyUserType{}, user)
		next.ServeHTTP(w, r.WithContext(ctx)) // token check successed
	})
}

// GetUserFromAuthorization returns error if authentication failed
func (h *Handler) getUserFromAuthorization(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(contextKeyUserType{}).(*models.User)
	if !ok {
		h.logger.Error("no User")
		return nil, errors.New("failed to authenticate")
	}
	if user == nil {
		h.logger.Error("nil User")
		return nil, errors.New("failed to authenticate")
	}

	return user, nil
}
