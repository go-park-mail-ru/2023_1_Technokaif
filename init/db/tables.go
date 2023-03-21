package db

type Tables interface {
	Users() string
	Tracks() string
	Artists() string
	Albums() string
	Listens() string
	ArtistsAlbums() string
	ArtistsTracks() string
	LikedAlbums() string
	LikedArtists() string
	LikedTracks() string
}
