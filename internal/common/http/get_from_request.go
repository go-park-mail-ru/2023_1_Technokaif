package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

const (
	TrackIdUrlParam  = "trackID"
	ArtistIdUrlParam = "artistID"
	AlbumIdUrlParam  = "albumID"
	UserIdUrlParam   = "userID"
)

// GetUserFromAuthorization returns error if authentication failed
func GetUserFromRequest(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(models.ContextKeyUserType{}).(*models.User)
	if !ok {
		return nil, errors.New("(middleware) no User in context")
	}
	if user == nil {
		return nil, errors.New("(middleware) nil User in context")
	}

	return user, nil
}

func GetTrackIDFromRequest(r *http.Request) (uint32, error) {
	return convertID(chi.URLParam(r, TrackIdUrlParam))
}

func GetUserIDFromRequest(r *http.Request) (uint32, error) {
	return convertID(chi.URLParam(r, UserIdUrlParam))
}

func GetArtistIDFromRequest(r *http.Request) (uint32, error) {
	return convertID(chi.URLParam(r, ArtistIdUrlParam))
}

func GetAlbumIDFromRequest(r *http.Request) (uint32, error) {
	return convertID(chi.URLParam(r, AlbumIdUrlParam))
}

func convertID(idUrl string) (uint32, error) {
	id, err := strconv.Atoi(idUrl)

	if err != nil {
		return 0, errors.New("invalid ID url param")
	}
	if id <= 0 {
		return 0, errors.New("invalid ID url param")
	}

	return uint32(id), nil
}
