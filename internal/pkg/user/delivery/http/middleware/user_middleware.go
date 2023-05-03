package middleware

import (
	"net/http"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Middleware struct {
	logger logger.Logger
}

func NewMiddleware(l logger.Logger) *Middleware {
	return &Middleware{
		logger: l,
	}
}

func (m *Middleware) CheckUserAuthAndResponce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlID, err := commonHTTP.GetUserIDFromRequest(r)
		if err != nil {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				commonHTTP.InvalidURLParameter, http.StatusBadRequest, m.logger, err)
			return
		}

		user, err := commonHTTP.GetUserFromRequest(r)
		if err != nil {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				commonHTTP.UnathorizedUser, http.StatusUnauthorized, m.logger, err)
			return
		}

		if urlID != user.ID {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				commonHTTP.ForbiddenUser, http.StatusForbidden, m.logger, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
