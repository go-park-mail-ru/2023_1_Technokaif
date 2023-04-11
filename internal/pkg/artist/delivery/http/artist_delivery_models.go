package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Create
type artistCreateInput struct {
	Name      string `json:"name" valid:"required"`
	AvatarSrc string `json:"cover" valid:"required"`
}

func (a *artistCreateInput) validate() error {
	a.Name = html.EscapeString(a.Name)
	a.AvatarSrc = html.EscapeString(a.AvatarSrc)

	_, err := valid.ValidateStruct(a)

	return err
}

func (aci *artistCreateInput) ToArtist(userID *uint32) models.Artist {
	return models.Artist{
		UserID:    userID,
		Name:      aci.Name,
		AvatarSrc: aci.AvatarSrc,
	}
}

type artistCreateResponse struct {
	ID uint32 `json:"id"`
}

// Delete
type artistDeleteResponse struct {
	Status string `json:"status"`
}

// Likes
type artistLikeResponse struct {
	Status string `json:"status"`
}
