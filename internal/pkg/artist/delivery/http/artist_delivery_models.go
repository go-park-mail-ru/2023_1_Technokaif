package http

import (
	"context"
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate easyjson -no_std_marshalers artist_delivery_models.go

// Response messages
const (
	albumNotFound  = "no such album"
	artistNotFound = "no such artist"
	trackNotFound  = "no such track"

	artistCreateNorights = "no rights to create artist"
	artistDeleteNoRights = "no rights to delete artist"

	artistCreateServerError = "can't create artist"
	artistGetServerError    = "can't get artist"
	artistsGetServerError   = "can't get artists"
	artistDeleteServerError = "can't delete artist"

	artistDeletedSuccessfully = "ok"
)

//easyjson:json
type ArtistTransfer struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	IsLiked   bool   `json:"isLiked"`
	AvatarSrc string `json:"cover"`
}

//easyjson:json
type ArtistTransfers []ArtistTransfer

type ArtistLikeChecker func(ctx context.Context, artistID, userID uint32) (bool, error)

// ArtistTransferFromEntry converts models.Artist to ArtistTransfer
func ArtistTransferFromEntry(ctx context.Context, a models.Artist, user *models.User,
	likeChecker ArtistLikeChecker) (ArtistTransfer, error) {

	var isLiked bool
	var err error

	if user != nil {
		isLiked, err = likeChecker(ctx, a.ID, user.ID)
		if err != nil {
			return ArtistTransfer{}, err
		}
	}

	return ArtistTransfer{
		ID:        a.ID,
		Name:      a.Name,
		IsLiked:   isLiked,
		AvatarSrc: a.AvatarSrc,
	}, nil
}

// ArtistTransferFromList converts []models.Artist to []ArtistTransfer
func ArtistTransferFromList(ctx context.Context, artists []models.Artist, user *models.User,
	likeChecker ArtistLikeChecker) (ArtistTransfers, error) {

	artistTransfers := make([]ArtistTransfer, 0, len(artists))
	for _, a := range artists {
		at, err := ArtistTransferFromEntry(ctx, a, user, likeChecker)
		if err != nil {
			return nil, err
		}

		artistTransfers = append(artistTransfers, at)
	}

	return artistTransfers, nil
}

//easyjson:json
type artistCreateInput struct {
	Name      string `json:"name" valid:"required"`
	AvatarSrc string `json:"cover" valid:"required"`
}

func (a *artistCreateInput) validateAndEscape() error {
	a.escapeHtml()

	_, err := valid.ValidateStruct(a)

	return err
}

func (a *artistCreateInput) escapeHtml() {
	a.Name = html.EscapeString(a.Name)
	a.AvatarSrc = html.EscapeString(a.AvatarSrc)
}

func (aci *artistCreateInput) ToArtist(userID *uint32) models.Artist {
	return models.Artist{
		UserID:    userID,
		Name:      aci.Name,
		AvatarSrc: aci.AvatarSrc,
	}
}

//easyjson:json
type artistCreateResponse struct {
	ID uint32 `json:"id"`
}

//easyjson:json
type artistDeleteResponse struct {
	Status string `json:"status"`
}

//easyjson:json
type artistLikeResponse struct {
	Status string `json:"status"`
}
