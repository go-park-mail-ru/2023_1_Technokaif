package middleware

import (
	"net/http"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
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
		user, err := commonHTTP.GetUserFromRequest(r)
		if err != nil {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				invalidAccessToken, http.StatusBadRequest, m.logger, err)
			return
		}

		csrfToken := r.Header.Get(csrfTokenHttpHeader)
		if csrfToken == "" {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				missingCSRFToken, http.StatusBadRequest, m.logger, err)
			return
		}

		userIDFromToken, err := m.tokenServices.CheckCSRFToken(csrfToken)
		if err != nil {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				invalidCSRFToken, http.StatusBadRequest, m.logger, err)
			return
		}
		if user.ID != userIDFromToken {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				invalidCSRFToken, http.StatusBadRequest, m.logger, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
