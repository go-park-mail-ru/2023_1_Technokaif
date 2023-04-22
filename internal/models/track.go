package models

type Track struct {
	ID            uint32  `db:"id"`
	Name          string  `db:"name"`
	AlbumID       *uint32 `db:"album_id"`
	AlbumPosition *uint32 `db:"album_position"`
	CoverSrc      string  `db:"cover_src"`
	RecordSrc     string  `db:"record_src"`
	Listens       uint32  `db:"listens"`
}

type TrackTransfer struct {
	ID            uint32           `json:"id"`
	Name          string           `json:"name"`
	AlbumID       *uint32          `json:"albumID,omitempty"`
	AlbumPosition *uint32          `json:"albumPosition,omitempty"`
	Artists       []ArtistTransfer `json:"artists"`
	CoverSrc      string           `json:"cover"`
	Listens       uint32           `json:"listens"`
	IsLiked       bool             `json:"isLiked"`
	RecordSrc     string           `json:"recordSrc"`
}

type artistsByTrackGetter func(trackID uint32) ([]Artist, error)
type trackLikeChecker func(trackID, userID uint32) (bool, error)

// TrackTransferFromEntry converts Track to TrackTransfer
func TrackTransferFromEntry(t Track, user *User, lc trackLikeChecker, ag artistsByTrackGetter) (TrackTransfer, error) {
	artists, err := ag(t.ID)
	if err != nil {
		return TrackTransfer{}, err
	}

	var isLiked = false
	if user != nil {
		isLiked, err = lc(t.ID, user.ID)
		if err != nil {
			return TrackTransfer{}, err
		}
	}

	return TrackTransfer{
		ID:            t.ID,
		Name:          t.Name,
		AlbumID:       t.AlbumID,
		AlbumPosition: t.AlbumPosition,
		Artists:       ArtistTransferFromQuery(artists),
		CoverSrc:      t.CoverSrc,
		Listens:       t.Listens,
		IsLiked:       isLiked,
		RecordSrc:     t.RecordSrc,
	}, nil
}

// TrackTransferFromQuery converts []Track to []TrackTransfer
func TrackTransferFromQuery(tracks []Track, user *User, lc trackLikeChecker, ag artistsByTrackGetter) ([]TrackTransfer, error) {
	trackTransfers := make([]TrackTransfer, 0, len(tracks))
	for _, t := range tracks {
		trackTransfer, err := TrackTransferFromEntry(t, user, lc, ag)
		if err != nil {
			return nil, err
		}

		trackTransfers = append(trackTransfers, trackTransfer)
	}

	return trackTransfers, nil
}
