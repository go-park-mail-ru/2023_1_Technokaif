package http

import (
	"net/http"
)

const AccessTokenCookieName = "X-ACCESS-Token"

func SetAccessTokenCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func GetAccessTokenFromCookie(r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil {
		return "", err
	}
	return tokenCookie.Value, nil
}
