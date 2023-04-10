package http

import (
	"net/http"
)

const acessTokenCookieName = "X-ACCESS-Token"

func SetAcessTokenCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     acessTokenCookieName,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func GetAcessTokenFromCookie(r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie(acessTokenCookieName)
	if err != nil {
		return "", err
	}
	return tokenCookie.Value, nil
}
