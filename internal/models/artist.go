package models

type Artist struct {
	ID        int
	Name      string
	AvatarSrc string
}

type ArtistFeed struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
