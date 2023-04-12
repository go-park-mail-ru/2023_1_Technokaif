package tests

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

func WrapRequestWithUserNotNilFunc(user *models.User) Wrapper {
	return func(req *http.Request) *http.Request {
		return WrapRequestWithUserNotNil(req, user)
	}
}

func WrapRequestWithUserFunc(user *models.User, doWrap bool) Wrapper {
	return func(req *http.Request) *http.Request {
		return WrapRequestWithUser(req, user, doWrap)
	}
}

func NoWrapUserFunc() Wrapper {
	return func(req *http.Request) *http.Request {
		return req
	}
}

func WrapRequestWithUser(r *http.Request, user *models.User, doWrap bool) *http.Request {
	if !doWrap {
		return r
	}
	return wrap(r, user)
}

func WrapRequestWithUserNotNil(r *http.Request, user *models.User) *http.Request {
	if user == nil {
		return r
	}
	return wrap(r, user)
}

func wrap(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
	return r.WithContext(ctx)
}
