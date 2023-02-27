package models

type Track struct {
	ID        uint
	Name      string
	AlbumID   string
	ArtistID  uint
	CoverSrc  string
	RecordSrc string
}

type TrackFeed struct {
	Name       string
	ArtistName string
}
