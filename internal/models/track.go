package models

type Track struct {
	ID        uint32 `db:"id"`
	Name      string `db:"name"`
	AlbumID   uint32 `db:"album_id"`
	CoverSrc  string `db:"cover_src"`
	RecordSrc string `db:"record_src"`
}

type TrackTransfer struct {
	ID        uint32           `json:"id"`
	Name      string           `json:"name"`
	AlbumID   uint32           `json:"albumID,omitempty"`
	Artists   []ArtistTransfer `json:"artists"`
	CoverSrc  string           `json:"cover"`
	RecordSrc string           `json:"record"`
}

type artistsByTrackGetter func(trackID uint32) ([]Artist, error)

// TrackTransferFromEntry converts Track to TrackTransfer
func TrackTransferFromEntry(t Track, artistsGetter artistsByTrackGetter) (TrackTransfer, error) {

	artists, err := artistsGetter(t.ID)
	if err != nil {
		return TrackTransfer{}, err
	}

	return TrackTransfer{
		ID:        t.ID,
		Name:      t.Name,
		AlbumID:   t.AlbumID,
		Artists:   ArtistTransferFromQuery(artists),
		CoverSrc:  t.CoverSrc,
		RecordSrc: t.RecordSrc,
	}, nil
}

// TrackTransferFromQuery converts []Track to []TrackTransfer
func TrackTransferFromQuery(tracks []Track, artistsGetter artistsByTrackGetter) ([]TrackTransfer, error) {
	trackTransfers := make([]TrackTransfer, 0, len(tracks))
	for _, t := range tracks {
		trackTransfer, err := TrackTransferFromEntry(t, artistsGetter)
		if err != nil {
			return nil, err
		}

		trackTransfers = append(trackTransfers, trackTransfer)
	}

	return trackTransfers, nil
}
