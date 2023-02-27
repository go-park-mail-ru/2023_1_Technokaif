package models

type Album struct {
	ID       int
	Name     string
	Info     string
	CoverSrc string
}

type AlbumFeed struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Artists     []ArtistFeed `json:"artists"`
	Description string       `json:"description"`
}
