package http

import (
	"errors"
	"html"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate easyjson -no_std_marshalers track_delivery_models.go

// Response messages
const (
	albumNotFound    = "no such album"
	artistNotFound   = "no such artist"
	playlistNotFound = "no such playlist"
	trackNotFound    = "no such track"

	trackCreateNorights = "no rights to create track"
	trackDeleteNoRights = "no rights to delete track"

	trackCreateServerError = "can't create track"
	trackGetServerError    = "can't get track"
	tracksGetServerError   = "can't get tracks"
	trackDeleteServerError = "can't delete track"
	listenAddError         = "can't add listen of track"

	trackDeletedSuccessfully = "ok"
	listenAddedSuccessfully  = "ok"
)

//easyjson:json
type trackCreateInput struct {
	Name          string   `json:"name" valid:"required"`
	AlbumID       *uint32  `json:"albumID"`
	AlbumPosition *uint32  `json:"albumPosition"`
	ArtistsID     []uint32 `json:"artistsID" valid:"required"`
	RecordSrc     string   `json:"record" valid:"required"`
}

//easyjson:json
type trackFeedInput struct {
	Days uint32 `json:"days" valid:"required"`
}

func (t *trackCreateInput) validateAndEscape() error {
	t.escapeHTML()

	if (t.AlbumID == nil) != (t.AlbumPosition == nil) {
		return errors.New("(delivery) albumID is nil while albumPosition isn't (or vice versa)")
	}

	_, err := valid.ValidateStruct(t)

	return err
}

func (t *trackCreateInput) escapeHTML() {
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

//easyjson:json
type trackCreateResponse struct {
	ID uint32 `json:"id"`
}

//easyjson:json
type trackDeleteResponse struct {
	Status string `json:"status"`
}

//easyjson:json
type trackLikeResponse struct {
	Status string `json:"status"`
}

//easyjson:json
type trackIncrementListensResponse struct {
	Status string `json:"status"`
}
