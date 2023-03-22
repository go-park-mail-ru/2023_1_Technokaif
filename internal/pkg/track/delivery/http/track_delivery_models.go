package http

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// Create
type trackCreateInput struct {
	Name      string   `json:"name"`
	AlbumID   uint32   `json:"albumID,omitempty"`
	ArtistsID []uint32 `json:"artistsID"`
	CoverSrc  string   `json:"cover"`
	RecordSrc string   `json:"record"`
}

func (tci *trackCreateInput) ToTrack() models.Track {
	return models.Track{
		Name:      tci.Name,
		AlbumID:   tci.AlbumID,
		CoverSrc:  tci.CoverSrc,
		RecordSrc: tci.RecordSrc,
	}
}

type trackCreateResponse struct {
	ID uint32 `json:"id"`
}

// Delete
type trackDeleteResponse struct {
	Status string `json:"status"`
}
