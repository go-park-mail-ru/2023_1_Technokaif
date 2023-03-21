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
	AlbumID   uint32           `json:"albumID,omitempty"` // TODO discuss
	Artists   []ArtistTransfer `json:"artists"`
	CoverSrc  string           `json:"cover"`
	RecordSrc string           `json:"record"`
}
