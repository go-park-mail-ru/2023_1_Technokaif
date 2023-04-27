package models

import (
	"context"
)

type Track struct {
	ID            uint32  `db:"id"`
	Name          string  `db:"name"`
	AlbumID       *uint32 `db:"album_id"`
	AlbumPosition *uint32 `db:"album_position"`
	CoverSrc      string  `db:"cover_src"`
	RecordSrc     string  `db:"record_src"`
	Duration      uint32  `db:"duration"`
	Listens       uint32  `db:"listens"`
}

type TrackTransfer struct {
	ID            uint32           `json:"id"`
	Name          string           `json:"name"`
	AlbumID       *uint32          `json:"albumID,omitempty"`
	AlbumPosition *uint32          `json:"albumPosition,omitempty"`
	Artists       []ArtistTransfer `json:"artists"`
	CoverSrc      string           `json:"cover"`
	Duration      uint32           `json:"duration"`
	Listens       uint32           `json:"listens"`
	IsLiked       bool             `json:"isLiked"`
	RecordSrc     string           `json:"recordSrc"`
}

type artistsByTrackGetter func(ctx context.Context, trackID uint32) ([]Artist, error)
type trackLikeChecker func(ctx context.Context, trackID, userID uint32) (bool, error)

// TrackTransferFromEntry converts Track to TrackTransfer
func TrackTransferFromEntry(ctx context.Context, t Track, user *User, likeChecker trackLikeChecker,
	artistLikeChecker artistLikeChecker, artistsGetter artistsByTrackGetter) (TrackTransfer, error) {

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

	at, err := ArtistTransferFromQuery(ctx, artists, user, artistLikeChecker)
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

// TrackTransferFromQuery converts []Track to []TrackTransfer
func TrackTransferFromQuery(ctx context.Context, tracks []Track, user *User, likeChecker trackLikeChecker,
	artistLikeChecker artistLikeChecker, artistsGetter artistsByTrackGetter) ([]TrackTransfer, error) {

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
