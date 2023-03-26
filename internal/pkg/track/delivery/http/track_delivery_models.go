package http

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Create
type trackCreateInput struct {
	Name      string   `json:"name" valid:"required"`
	AlbumID   uint32   `json:"albumID,omitempty"`
	ArtistsID []uint32 `json:"artistsID" valid:"required"`
	CoverSrc  string   `json:"cover" valid:"required"`
	RecordSrc string   `json:"record" valid:"required"`
}

func (t *trackCreateInput) validate() error {
	sanitizer := bluemonday.StrictPolicy()
	t.Name = sanitizer.Sanitize(t.Name)
	t.CoverSrc = sanitizer.Sanitize(t.CoverSrc)
	t.RecordSrc = sanitizer.Sanitize(t.RecordSrc)

	_, err := valid.ValidateStruct(t)

	return err
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
