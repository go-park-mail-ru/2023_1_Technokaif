package models

type Artist struct {
	ID        uint32
	Name      string
	AvatarSrc string
}

type ArtistTransfer struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	AvatarSrc string `json:"cover"`
}
