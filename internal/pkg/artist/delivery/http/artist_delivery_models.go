package http

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Create
type artistCreateInput struct {
	Name      string `json:"name" valid:"required"`
	AvatarSrc string `json:"cover" valid:"required"`
}

func (a *artistCreateInput) validate() error {
	sanitizer := bluemonday.StrictPolicy()
	a.Name = sanitizer.Sanitize(a.Name)
	a.AvatarSrc = sanitizer.Sanitize(a.AvatarSrc)

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
