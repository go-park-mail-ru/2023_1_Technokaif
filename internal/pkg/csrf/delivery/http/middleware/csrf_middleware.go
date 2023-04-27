package middleware

import (
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const csrfTokenHttpHeader = "X-CSRF-Token"

// Response messages
const (
	invalidAccessToken = "invalid access token"
	invalidCSRFToken   = "invalid CSRF token"

	missingCSRFToken = "missing CSRF token"
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
		if err != nil {
			commonHttp.ErrorResponseWithErrLogging(w, invalidAccessToken, http.StatusBadRequest, m.logger, err)
			return
		}

		csrfToken := r.Header.Get(csrfTokenHttpHeader)
		if csrfToken == "" {
			commonHttp.ErrorResponseWithErrLogging(w, missingCSRFToken, http.StatusBadRequest, m.logger, err)
			return
		}

		userIDFromToken, err := m.tokenServices.CheckCSRFToken(csrfToken)
		if err != nil {
			commonHttp.ErrorResponseWithErrLogging(w, invalidCSRFToken, http.StatusBadRequest, m.logger, err)
			return
		}
		if user.ID != userIDFromToken {
			m.logger.Info("invalid CSRF token: userID and userID from token are not equal")
			commonHttp.ErrorResponse(w, invalidCSRFToken, http.StatusBadRequest, m.logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}
