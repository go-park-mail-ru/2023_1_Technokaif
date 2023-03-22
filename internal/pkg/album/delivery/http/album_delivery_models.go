package http

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// Create
type albumCreateInput struct {
	Name        string   `json:"name"`
	ArtistsID   []uint32 `json:"artistsID"`
	Description string   `json:"description"`
	CoverSrc    string   `json:"cover"`
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
	ID          uint32   `json:"id"`
	Name        string   `json:"name"`
	ArtistsID   []uint32 `json:"artistsID"`
	Description string   `json:"description"`
	CoverSrc    string   `json:"cover"`
}

func (aci *albumChangeInput) ToAlbum() models.Album {
	return models.Album{
		ID:          aci.ID,
		Name:        aci.Name,
		Description: aci.Description,
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
