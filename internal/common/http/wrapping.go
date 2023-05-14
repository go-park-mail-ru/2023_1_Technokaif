package http

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type contextKeyReqIDType struct{}
type contextKeyUserType struct{}
type contextKeyResponseCodeType struct{}

func WrapUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), contextKeyUserType{}, user)
	return r.WithContext(ctx)
}

func WrapReqID(r *http.Request, id uint32) *http.Request {
	ctx := context.WithValue(r.Context(), contextKeyReqIDType{}, id)
	return r.WithContext(ctx)
}
