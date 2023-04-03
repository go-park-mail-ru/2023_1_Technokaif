package http

import (
	"net/http"
)

func SetAcessTokenCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "X-ACCESS-Token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api",
	}
	http.SetCookie(w, &cookie)
}
