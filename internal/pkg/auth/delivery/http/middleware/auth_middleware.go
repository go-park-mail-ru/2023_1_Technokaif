package middleware

import (
	"context"
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Middleware struct {
	authServices  auth.Usecase
	tokenServices token.Usecase
	logger        logger.Logger
}

func NewMiddleware(u auth.Usecase, t token.Usecase, l logger.Logger) *Middleware {
	return &Middleware{
		authServices:  u,
		tokenServices: t,
		logger:        l,
	}
}

// Authorization is HTTP middleware which sets a value on the request context
func (m *Middleware) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenCookie, err := r.Cookie(commonHttp.AcessTokenCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				m.logger.Infof("middleware: %s", err.Error())
				next.ServeHTTP(w, r) // no cookies
				return
			}

			m.logger.Errorf("middleware: %s", err.Error())
			commonHttp.ErrorResponse(w, "server error", http.StatusInternalServerError, m.logger) // server error
			return
		}
		if tokenCookie.Value == "" {
			m.logger.Infof("middleware: %s", "empty cookies")
			next.ServeHTTP(w, r) // empty cookies
			return
		}

		userId, userVersion, err := m.tokenServices.CheckAccessToken(tokenCookie.Value)
		if err != nil {
			m.logger.Infof("middleware: %s", err.Error())
			next.ServeHTTP(w, r) // token check failed
			return
		}

		user, err := m.authServices.GetUserByAuthData(userId, userVersion)
		if err != nil {
			var errNoSuchUser *models.NoSuchUserError
			if errors.As(err, &errNoSuchUser) {
				m.logger.Infof("middleware: %s", err.Error())
				next.ServeHTTP(w, r) // data check failed
				return
			}

			m.logger.Errorf("middleware: %s", err.Error())
			commonHttp.ErrorResponse(w, "server failed to check authorization", http.StatusInternalServerError, m.logger)
			return
		}

		m.logger.Infof("user version : %d", user.Version)

		ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
		next.ServeHTTP(w, r.WithContext(ctx)) // token check successed
	})
}
