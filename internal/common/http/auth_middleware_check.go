package common_http

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	authMiddleware "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/delivery/http/middleware"
)

// GetUserFromAuthorization returns error if authentication failed
func GetUserFromAuthorization(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(authMiddleware.ContextKeyUserType{}).(*models.User)
	if !ok {
		return nil, errors.New("(Middleware) no User in context")
	}
	if user == nil {
		return nil, errors.New("(Middleware) nil User in context")
	}

	return user, nil
}
