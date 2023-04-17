package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Create
type playlistCreateInput struct {
	Name        string   `json:"name" valid:"required"`
	UsersID     []uint32 `json:"users" valid:"required"`
	Description *string  `json:"description"`
	CoverSrc    *string  `json:"cover"`
}

func (pci *playlistCreateInput) validate() error {
	pci.Name = html.EscapeString(pci.Name)
	if pci.Description != nil {
		*pci.Description = html.EscapeString(*pci.Description)
	}
	if pci.CoverSrc != nil {
		*pci.CoverSrc = html.EscapeString(*pci.CoverSrc)
	}

	_, err := valid.ValidateStruct(pci)

	return err
}

func (pci *playlistCreateInput) ToPlaylist() models.Playlist {
	return models.Playlist{
		Name:        pci.Name,
		Description: pci.Description,
		CoverSrc:    pci.CoverSrc,
	}
}

// Update
type playlistUpdateInput struct {
	ID          uint32   `json:"id" valid:"required"`
	Name        string   `json:"name" valid:"required"`
	UsersID     []uint32 `json:"users" valid:"required"`
	Description *string  `json:"description"`
	CoverSrc    *string  `json:"cover"`
}

func (pui *playlistUpdateInput) validate() error {
	pui.Name = html.EscapeString(pui.Name)
	if pui.Description != nil {
		*pui.Description = html.EscapeString(*pui.Description)
	}
	if pui.CoverSrc != nil {
		*pui.CoverSrc = html.EscapeString(*pui.CoverSrc)
	}

	_, err := valid.ValidateStruct(pui)

	return err
}

func (pui *playlistUpdateInput) ToPlaylist() models.Playlist {
	return models.Playlist{
		ID:          pui.ID,
		Name:        pui.Name,
		Description: pui.Description,
		CoverSrc:    pui.CoverSrc,
	}
}

type playlistCreateResponse struct {
	ID uint32 `json:"id"`
}

type defaultResponse struct {
	Status string `json:"status"`
}
