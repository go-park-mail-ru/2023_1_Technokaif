package middleware

import (
	"errors"
	"net/http"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Response messages
const (
	tokenCheckFail    = "token check failed"
	authDataCheckFail = "auth data check failed"

	tokenGetServerError  = "server can't get access token"
	authCheckServerErorr = "server can't check authorization"
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

		logger: l,
	}
}

// Authorization is HTTP middleware which sets a value on the request context
func (m *Middleware) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := commonHTTP.GetAccessTokenFromCookie(r)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				m.logger.Infof("middleware: %v", err)
				next.ServeHTTP(w, r) // no cookies
				return
			}

			m.logger.Errorf("middleware: %v", err)
			commonHTTP.ErrorResponse(w, r, tokenGetServerError, http.StatusInternalServerError, m.logger)
			return
		}
		if token == "" {
			m.logger.Infof("middleware: %s", "empty cookies")
			next.ServeHTTP(w, r) // empty cookies
			return
		}

		userId, userVersion, err := m.tokenServices.CheckAccessToken(token)
		if err != nil {
			m.logger.Infof("middleware: %v", err)
			commonHTTP.SetAccessTokenCookie(w, "")
			commonHTTP.ErrorResponse(w, r, tokenCheckFail, http.StatusBadRequest, m.logger) // token check failed
			return
		}

		user, err := m.authServices.GetUserByAuthData(r.Context(), userId, userVersion)
		if err != nil {
			var errNoSuchUser *models.NoSuchUserError
			if errors.As(err, &errNoSuchUser) {
				m.logger.Infof("middleware: %v", err)
				commonHTTP.SetAccessTokenCookie(w, "")
				commonHTTP.ErrorResponse(w, r, authDataCheckFail, http.StatusBadRequest, m.logger) // auth data check failed
				return
			}

			m.logger.Errorf("middleware: %v", err)
			commonHTTP.ErrorResponse(w, r, authCheckServerErorr, http.StatusInternalServerError, m.logger)
			return
		}

		m.logger.InfofReqID(r.Context(), "user version : %d", user.Version)

		reqWithUser := commonHTTP.WrapUser(r, user)
		next.ServeHTTP(w, commonHTTP.WrapWithCode(reqWithUser, commonHTTP.GetResponseCodeFromRequest(r))) // token check successed
	})
}
