package http

import (
	"context"
	"errors"
	"html"

	valid "github.com/asaskevich/govalidator"

	artistHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/delivery/http"

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

	trackDeletedSuccessfully = "ok"
)

//easyjson:json
type TrackTransfer struct {
	ID            uint32                     `json:"id"`
	Name          string                     `json:"name"`
	AlbumID       *uint32                    `json:"albumID,omitempty"`
	AlbumPosition *uint32                    `json:"albumPosition,omitempty"`
	Artists       artistHTTP.ArtistTransfers `json:"artists"`
	CoverSrc      string                     `json:"cover"`
	Duration      uint32                     `json:"duration"`
	Listens       uint32                     `json:"listens"`
	IsLiked       bool                       `json:"isLiked"`
	RecordSrc     string                     `json:"recordSrc"`
}

//easyjson:json
type TrackTransfers []TrackTransfer

type artistsByTrackGetter func(ctx context.Context, trackID uint32) ([]models.Artist, error)
type trackLikeChecker func(ctx context.Context, trackID, userID uint32) (bool, error)

// TrackTransferFromEntry converts Track to TrackTransfer
func TrackTransferFromEntry(ctx context.Context, t models.Track, user *models.User, likeChecker trackLikeChecker,
	artistLikeChecker artistHTTP.ArtistLikeChecker, artistsGetter artistsByTrackGetter) (TrackTransfer, error) {

	artists, err := artistsGetter(ctx, t.ID)
	if err != nil {
		return TrackTransfer{}, err
	}

	var isLiked = false
	if user != nil {
		isLiked, err = likeChecker(ctx, t.ID, user.ID)
		if err != nil {
			return TrackTransfer{}, err
		}
	}

	at, err := artistHTTP.ArtistTransferFromList(ctx, artists, user, artistLikeChecker)
	if err != nil {
		return TrackTransfer{}, err
	}

	return TrackTransfer{
		ID:            t.ID,
		Name:          t.Name,
		AlbumID:       t.AlbumID,
		AlbumPosition: t.AlbumPosition,
		Artists:       at,
		CoverSrc:      t.CoverSrc,
		Duration:      t.Duration,
		Listens:       t.Listens,
		IsLiked:       isLiked,
		RecordSrc:     t.RecordSrc,
	}, nil
}

// TrackTransferFromList converts []Track to []TrackTransfer
func TrackTransferFromList(ctx context.Context, tracks []models.Track, user *models.User, likeChecker trackLikeChecker,
	artistLikeChecker artistHTTP.ArtistLikeChecker, artistsGetter artistsByTrackGetter) (TrackTransfers, error) {

	trackTransfers := make([]TrackTransfer, 0, len(tracks))
	for _, t := range tracks {
		trackTransfer, err := TrackTransferFromEntry(ctx, t, user, likeChecker, artistLikeChecker, artistsGetter)
		if err != nil {
			return nil, err
		}

		trackTransfers = append(trackTransfers, trackTransfer)
	}

	return trackTransfers, nil
}

//easyjson:json
type trackCreateInput struct {
	Name          string   `json:"name" valid:"required"`
	AlbumID       *uint32  `json:"albumID"`
	AlbumPosition *uint32  `json:"albumPosition"`
	ArtistsID     []uint32 `json:"artistsID" valid:"required"`
	RecordSrc     string   `json:"record" valid:"required"`
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
