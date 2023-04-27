package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

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

// Create
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

type albumCreateResponse struct {
	ID uint32 `json:"id"`
}

// Delete
type albumDeleteResponse struct {
	Status string `json:"status"`
}

// Likes
type albumLikeResponse struct {
	Status string `json:"status"`
}
