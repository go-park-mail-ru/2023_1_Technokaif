package models

type Track struct {
	ID        uint32
	Name      string
	AlbumID   uint32
	CoverSrc  string
	RecordSrc string
}

type TrackTransfer struct {
	ID        uint32           `json:"id"`
	Name      string           `json:"name"`
	Artists   []ArtistTransfer `json:"artists"`
	CoverSrc  string           `json:"cover"`
	RecordSrc string           `json:"record"`
}
