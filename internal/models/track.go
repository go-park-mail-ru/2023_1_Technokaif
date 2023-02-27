package models

type Track struct {
	ID          uint
	Name        string
	ArtistID    string
	CoverSource string
}

type TrackFeed struct {
	Name       string
	ArtistName string
}
