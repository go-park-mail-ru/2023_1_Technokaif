package models

import "context"

//go:generate easyjson -no_std_marshalers track.go

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

//easyjson:json
type TrackTransfer struct {
	ID            uint32          `json:"id"`
	Name          string          `json:"name"`
	AlbumID       *uint32         `json:"albumID,omitempty"`
	AlbumPosition *uint32         `json:"albumPosition,omitempty"`
	Artists       ArtistTransfers `json:"artists"`
	CoverSrc      string          `json:"cover"`
	Duration      uint32          `json:"duration"`
	Listens       uint32          `json:"listens"`
	IsLiked       bool            `json:"isLiked"`
	RecordSrc     string          `json:"recordSrc"`
}

//easyjson:json
type TrackTransfers []TrackTransfer

type artistsByTrackGetter func(ctx context.Context, trackID uint32) ([]Artist, error)
type trackLikeChecker func(ctx context.Context, trackID, userID uint32) (bool, error)

// TrackTransferFromEntry converts Track to TrackTransfer
func TrackTransferFromEntry(ctx context.Context, t Track, user *User, likeChecker trackLikeChecker,
	artistLikeChecker ArtistLikeChecker, artistsGetter artistsByTrackGetter) (TrackTransfer, error) {

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

	at, err := ArtistTransferFromList(ctx, artists, user, artistLikeChecker)
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
func TrackTransferFromList(ctx context.Context, tracks []Track, user *User, likeChecker trackLikeChecker,
	artistLikeChecker ArtistLikeChecker, artistsGetter artistsByTrackGetter) (TrackTransfers, error) {

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
