package models

type Album struct {
	ID          int
	Name        string
	Description string
	CoverSrc    string
}

type AlbumTransfer struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Artists     []ArtistTransfer `json:"artists"`
	Description string           `json:"description"`
	CoverSrc    string           `json:"cover"`
}
