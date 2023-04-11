package http

import (
	"net/http"
)

const AcessTokenCookieName = "X-ACCESS-Token"

func SetAcessTokenCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     AcessTokenCookieName,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func GetAcessTokenFromCookie(r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie(AcessTokenCookieName)
	if err != nil {
		return "", err
	}
	return tokenCookie.Value, nil
}
