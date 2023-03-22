package http

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// Create
type artistCreateInput struct {
	Name      string `json:"name"`
	AvatarSrc string `json:"avatar"`
}

func (aci *artistCreateInput) ToArtist() models.Artist {
	return models.Artist{
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
