package models

type Track struct {
	ID        int
	Name      string
	AlbumID   string
	ArtistID  int
	CoverSrc  string
	RecordSrc string
}

type TrackFeed struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Artists  []ArtistFeed `json:"artists"`
	CoverSrc string       `json:"cover"`
}
