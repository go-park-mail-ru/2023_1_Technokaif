package models

type Artist struct {
	ID        uint32 `db:"id"`
	Name      string `db:"name"`
	AvatarSrc string `db:"avatar_src"`
}

type ArtistTransfer struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	AvatarSrc string `json:"cover"`
}
