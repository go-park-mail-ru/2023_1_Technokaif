package http

import (
	"errors"
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Create
type trackCreateInput struct {
	Name          string   `json:"name" valid:"required"`
	AlbumID       *uint32  `json:"albumID"`
	AlbumPosition *uint32  `json:"albumPosition"`
	ArtistsID     []uint32 `json:"artistsID" valid:"required"`
	RecordSrc     string   `json:"record" valid:"required"`
}

func (t *trackCreateInput) validateAndEscape() error {
	t.escapeHtml()

	if (t.AlbumID == nil) != (t.AlbumPosition == nil) {
		return errors.New("(delivery) albumID is nil while albumPosition isn't (or vice versa)")
	}

	_, err := valid.ValidateStruct(t)

	return err
}

func (t *trackCreateInput) escapeHtml() {
	t.Name = html.EscapeString(t.Name)
}


func (tci *trackCreateInput) ToTrack() models.Track {
	return models.Track{
		Name:          tci.Name,
		AlbumID:       tci.AlbumID,
		AlbumPosition: tci.AlbumPosition,
		RecordSrc:     tci.RecordSrc,
	}
}

type trackCreateResponse struct {
	ID uint32 `json:"id"`
}

// Delete
type trackDeleteResponse struct {
	Status string `json:"status"`
}

// Likes
type trackLikeResponse struct {
	Status string `json:"status"`
}
