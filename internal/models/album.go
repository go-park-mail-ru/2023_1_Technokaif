package models

type Album struct {
	ID          uint32 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CoverSrc    string `db:"cover_src"`
}

type AlbumTransfer struct {
	ID          uint32           `json:"id"`
	Name        string           `json:"name"`
	Artists     []ArtistTransfer `json:"artists"`
	Description string           `json:"description"`
	CoverSrc    string           `json:"cover"`
}
