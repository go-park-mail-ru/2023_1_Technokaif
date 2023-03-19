package auth_middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
)

type AuthMiddleware struct {
	authServices auth.AuthUsecase
	logger   	 logger.Logger
}

func NewAuthMiddleware(u auth.AuthUsecase, l logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authServices: u,
		logger:   l,
	}
}

type ContextKeyUserType struct{}

// Authorization is HTTP middleware which sets a value on the request context
func (m *AuthMiddleware) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer"
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)
		reqToken = strings.ReplaceAll(reqToken, " ", "")

		m.logger.Info("auth token : " + reqToken)

		if authHeader == "" || reqToken == authHeader || reqToken == "" {
			m.logger.Info("middleware: missing token")
			next.ServeHTTP(w, r) // missing token
			return
		}

		userId, userVersion, err := m.authServices.CheckAccessToken(reqToken)
		if err != nil {
			m.logger.Infof("middleware: %s", err.Error())
			next.ServeHTTP(w, r) // token check failed
			return
		}

		user, err := m.authServices.GetUserByAuthData(userId, userVersion)
		if err != nil {
			m.logger.Infof("middleware: %s", err.Error())
			next.ServeHTTP(w, r) // UserAuth data check failed
			return
		}

		m.logger.Infof("user version : %d", user.Version)

		ctx := context.WithValue(r.Context(), ContextKeyUserType{}, user)
		next.ServeHTTP(w, r.WithContext(ctx)) // token check successed
	})
}
