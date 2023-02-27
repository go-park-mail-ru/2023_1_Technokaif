package models

type Album struct {
	ID       uint
	Name     string
	Info     string
	CoverSrc string
}

type AlbumFeed struct {
	Name       string
	ArtistName string
}
