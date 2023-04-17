package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Create
type playlistCreateInput struct {
	Name        string   `json:"name" valid:"required"`
	UsersID     []uint32 `json:"usersID" valid:"required"`
	Description *string  `json:"description"`
	CoverSrc    *string  `json:"cover" valid:"required"`
}

func (p *playlistCreateInput) validate() error {
	p.Name = html.EscapeString(p.Name)
	if p.Description != nil {
		*p.Description = html.EscapeString(*p.Description)
	}
	if p.CoverSrc != nil {
		*p.CoverSrc = html.EscapeString(*p.CoverSrc)
	}

	_, err := valid.ValidateStruct(p)

	return err
}

func (pci *playlistCreateInput) ToPlaylist() models.Playlist {
	return models.Playlist{
		Name:        pci.Name,
		Description: pci.Description,
		CoverSrc:    pci.CoverSrc,
	}
}

type playlistCreateResponse struct {
	ID uint32 `json:"id"`
}

type defaultResponse struct {
	Status string `json:"status"`
}
