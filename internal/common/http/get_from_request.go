package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

const (
	TrackIdUrlParam    = "trackID"
	ArtistIdUrlParam   = "artistID"
	AlbumIdUrlParam    = "albumID"
	PlaylistIdUrlParam = "playlistID"
	UserIdUrlParam     = "userID"
)

var ErrUnauthorized = &models.UnathorizedError{}

// GetUserFromRequest returns error if authentication failed
func GetUserFromRequest(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(contextKeyUserType{}).(*models.User)
	if !ok {
		return nil, ErrUnauthorized
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	return user, nil
}

func GetReqIDFromContext(ctx context.Context) (uint32, error) {
	reqID, ok := ctx.Value(contextKeyReqIDType{}).(uint32)
	if !ok {
		return 0, errors.New("no reqID in context")
	}

	return reqID, nil
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

func GetPlaylistIDFromRequest(r *http.Request) (uint32, error) {
	return convertID(chi.URLParam(r, PlaylistIdUrlParam))
}

func convertID(idUrl string) (uint32, error) {
	id, err := strconv.ParseUint(idUrl, 10, 32)
	if err != nil || id == 0 {
		return 0, errors.New("invalid ID url param")
	}

	return uint32(id), nil
}
