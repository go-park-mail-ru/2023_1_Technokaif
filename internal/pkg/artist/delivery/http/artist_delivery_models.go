package http

import (
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
