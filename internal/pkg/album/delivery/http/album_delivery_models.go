package http

import (
	"context"
	"html"

	valid "github.com/asaskevich/govalidator"
	artistHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate easyjson -no_std_marshalers album_delivery_models.go

// Response messages
const (
	albumNotFound  = "no such album"
	artistNotFound = "no such artist"

	albumCreateNorights = "no rights to create album"
	albumDeleteNoRights = "no rights to delete album"

	albumCreateServerError = "can't create album"
	albumGetServerError    = "can't get album"
	albumsGetServerError   = "can't get albums"
	albumDeleteServerError = "can't delete album"

	albumDeletedSuccessfully = "ok"
)

//easyjson:json
type AlbumTransfer struct {
	ID          uint32                     `json:"id"`
	Name        string                     `json:"name"`
	Artists     artistHTTP.ArtistTransfers `json:"artists"`
	Description *string                    `json:"description,omitempty"`
	IsLiked     bool                       `json:"isLiked"`
	CoverSrc    string                     `json:"cover"`
}

//easyjson:json
type AlbumTransfers []AlbumTransfer

type artistsByAlbumGetter func(ctx context.Context, albumID uint32) ([]models.Artist, error)
type albumLikeChecker func(ctx context.Context, albumID, userID uint32) (bool, error)

// AlbumTransferFromEntry converts models.Album to AlbumTransfer
func AlbumTransferFromEntry(ctx context.Context, a models.Album, user *models.User, likeChecker albumLikeChecker,
	artistLikeChecker artistHTTP.ArtistLikeChecker, artistsGetter artistsByAlbumGetter) (AlbumTransfer, error) {

	artists, err := artistsGetter(ctx, a.ID)
	if err != nil {
		return AlbumTransfer{}, err
	}

	var isLiked = false
	if user != nil {
		isLiked, err = likeChecker(ctx, a.ID, user.ID)
		if err != nil {
			return AlbumTransfer{}, err
		}
	}

	at, err := artistHTTP.ArtistTransferFromList(ctx, artists, user, artistLikeChecker)
	if err != nil {
		return AlbumTransfer{}, err
	}

	return AlbumTransfer{
		ID:          a.ID,
		Name:        a.Name,
		Artists:     at,
		Description: a.Description,
		IsLiked:     isLiked,
		CoverSrc:    a.CoverSrc,
	}, nil
}

// AlbumTransferFromList converts []models.Album to []AlbumTransfer
func AlbumTransferFromList(ctx context.Context, albums []models.Album, user *models.User, likeChecker albumLikeChecker,
	artistLikeChecker artistHTTP.ArtistLikeChecker, artistsGetter artistsByAlbumGetter) (AlbumTransfers, error) {

	albumTransfers := make([]AlbumTransfer, 0, len(albums))
	for _, a := range albums {
		at, err := AlbumTransferFromEntry(ctx, a, user, likeChecker, artistLikeChecker, artistsGetter)
		if err != nil {
			return nil, err
		}

		albumTransfers = append(albumTransfers, at)
	}

	return albumTransfers, nil
}

//easyjson:json
type albumCreateInput struct {
	Name        string   `json:"name" valid:"required"`
	ArtistsID   []uint32 `json:"artists" valid:"required"`
	Description *string  `json:"description"`
	CoverSrc    string   `json:"cover" valid:"required"`
}

func (a *albumCreateInput) validateAndEscape() error {
	a.escapeHtml()

	_, err := valid.ValidateStruct(a)

	return err
}

func (a *albumCreateInput) escapeHtml() {
	a.Name = html.EscapeString(a.Name)
	if a.Description != nil {
		*a.Description = html.EscapeString(*a.Description)
	}
	a.CoverSrc = html.EscapeString(a.CoverSrc)
}

func (aci *albumCreateInput) ToAlbum() models.Album {
	return models.Album{
		Name:        aci.Name,
		Description: aci.Description,
		CoverSrc:    aci.CoverSrc,
	}
}

//easyjson:json
type albumCreateResponse struct {
	ID uint32 `json:"id"`
}

//easyjson:json
type albumDeleteResponse struct {
	Status string `json:"status"`
}

//easyjson:json
type albumLikeResponse struct {
	Status string `json:"status"`
}
