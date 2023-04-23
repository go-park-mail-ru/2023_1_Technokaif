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
			m.logger.Infof("invalid url parameter: %v", err.Error())
			commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, m.logger)
			return
		}

		user, err := commonHttp.GetUserFromRequest(r)
		if err != nil {
			m.logger.Infof("unathorized user: %v", err)
			commonHttp.ErrorResponse(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, m.logger)
			return
		}

		if urlID != user.ID {
			m.logger.Infof("forbidden user with id #%d", urlID)
			commonHttp.ErrorResponse(w, "user has no rights", http.StatusForbidden, m.logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}
