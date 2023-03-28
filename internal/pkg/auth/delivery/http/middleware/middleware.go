package middleware

import (
	"context"
	"net/http"
	"strings"
	"errors"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Middleware struct {
	authServices auth.Usecase
	logger       logger.Logger
}

func NewMiddleware(u auth.Usecase, l logger.Logger) *Middleware {
	return &Middleware{
		authServices: u,
		logger:       l,
	}
}

// Authorization is HTTP middleware which sets a value on the request context
func (m *Middleware) Authorization(next http.Handler) http.Handler {
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
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			m.logger.Infof("middleware: %s", err.Error())
			next.ServeHTTP(w, r) // UserAuth data check failed
			return
		} else if err != nil {
			m.logger.Errorf("middleware: %s", err.Error())
			commonHttp.ErrorResponse(w, "server failed to check authorization", http.StatusInternalServerError, m.logger)
			return
		}

		m.logger.Infof("user version : %d", user.Version)

		ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
		next.ServeHTTP(w, r.WithContext(ctx)) // token check successed
	})
}
