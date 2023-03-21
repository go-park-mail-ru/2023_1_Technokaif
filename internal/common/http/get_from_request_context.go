package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

const (
	trackIdUrlParam  = "trackID"
	artistIdUrlParam = "artistID"
	albumIdUrlParam  = "albumID"
	userIdUrlParam   = "userID"
)

// GetUserFromAuthorization returns error if authentication failed
func GetUserFromRequest(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(models.ContextKeyUserType{}).(*models.User)
	if !ok {
		return nil, errors.New("(Middleware) no User in context")
	}
	if user == nil {
		return nil, errors.New("(Middleware) nil User in context")
	}

	return user, nil
}

func GetTrackIDFromRequest(r *http.Request) (uint32, error) {
	return convertID(chi.URLParam(r, trackIdUrlParam))
}

func GetUserIDFromRequest(r *http.Request) (uint32, error) {
	return convertID(chi.URLParam(r, userIdUrlParam))
}

func GetArtistIDFromRequest(r *http.Request) (uint32, error) {
	return convertID(chi.URLParam(r, artistIdUrlParam))
}

func convertID(idUrl string) (uint32, error) {
	id, err := strconv.Atoi(idUrl)

	if id <= 0 {
		return 0, errors.New("invalid ID url param")
	}
	if err != nil {
		return 0, errors.New("invalid ID url param")
	}

	return uint32(id), nil
}
