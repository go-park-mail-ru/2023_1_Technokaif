package middleware

import (
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Middleware struct {
	tokenServices token.Usecase
	logger        logger.Logger
}

func NewMiddleware(t token.Usecase, l logger.Logger) *Middleware {
	return &Middleware{
		tokenServices: t,
		logger:        l,
	}
}

func (m *Middleware) CheckCSRFToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := commonHttp.GetUserFromRequest(r)

		csrfToken := r.Header.Get("X-CSRF-Token")

		userIDFromToken, err := m.tokenServices.CheckCSRFToken(csrfToken); 
		if err != nil {
			commonHttp.ErrorResponseWithErrLogging(w, "invalid CSRF token", http.StatusBadRequest, m.logger, err)
			return
		}
		if user.ID != userIDFromToken {
			m.logger.Info("invalid CSRF token: userID and userID from token are not equal")
			commonHttp.ErrorResponse(w, "invalid CSRF token", http.StatusBadRequest, m.logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}
