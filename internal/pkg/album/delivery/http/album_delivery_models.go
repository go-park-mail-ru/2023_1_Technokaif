package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Create
type albumCreateInput struct {
	Name        string   `json:"name" valid:"required"`
	ArtistsID   []uint32 `json:"artistsID" valid:"required"`
	Description *string  `json:"description"`
	CoverSrc    string   `json:"cover" valid:"required"`
}

func (a *albumCreateInput) validate() error {
	a.Name = html.EscapeString(a.Name)
	if a.Description != nil {
		*a.Description = html.EscapeString(*a.Description)
	}
	a.CoverSrc = html.EscapeString(a.CoverSrc)

	_, err := valid.ValidateStruct(a)

	return err
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

// Change
type albumChangeInput struct {
	ID          uint32   `json:"id" valid:"required"`
	Name        string   `json:"name" valid:"required"`
	ArtistsID   []uint32 `json:"artistsID" valid:"required"`
	Description string   `json:"description,omitempty"`
	CoverSrc    string   `json:"cover" valid:"required"`
}

func (a *albumChangeInput) validate() error {
	a.Name = html.EscapeString(a.Name)
	a.Description = html.EscapeString(a.Description)
	a.CoverSrc = html.EscapeString(a.CoverSrc)

	_, err := valid.ValidateStruct(a)

	return err
}

func (aci *albumChangeInput) ToAlbum() models.Album {
	return models.Album{
		ID:          aci.ID,
		Name:        aci.Name,
		Description: &aci.Description,
		CoverSrc:    aci.CoverSrc,
	}
}

type albumChangeResponse struct {
	Message string `json:"status"`
}

// Delete
type albumDeleteResponse struct {
	Status string `json:"status"`
}

// Likes
type albumLikeResponse struct {
	Status string `json:"status"`
}
