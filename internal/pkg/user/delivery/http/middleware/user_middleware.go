package middleware

import (
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
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
		urlID, err := commonHttp.GetUserIDFromRequest(r)
		if err != nil {
			commonHttp.ErrorResponseWithErrLogging(w, r,
				commonHttp.InvalidURLParameter, http.StatusBadRequest, m.logger, err)
			return
		}

		user, err := commonHttp.GetUserFromRequest(r)
		if err != nil {
			commonHttp.ErrorResponseWithErrLogging(w, r,
				commonHttp.UnathorizedUser, http.StatusUnauthorized, m.logger, err)
			return
		}

		if urlID != user.ID {
			commonHttp.ErrorResponseWithErrLogging(w, r,
				commonHttp.ForbiddenUser, http.StatusForbidden, m.logger, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
