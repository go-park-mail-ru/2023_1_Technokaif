package track

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

//go:generate mockgen -source=track.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with tracks
type Usecase interface {
	Create(track models.Track, artistsID []uint32, userID uint32) (uint32, error)
	GetByID(trackID uint32) (*models.Track, error)
	Delete(trackID uint32, userID uint32) error
	GetFeed() ([]models.Track, error)
	GetByAlbum(albumID uint32) ([]models.Track, error)
	GetByPlaylist(playlistID uint32) ([]models.Track, error)
	GetByArtist(artistID uint32) ([]models.Track, error)
	GetLikedByUser(userID uint32) ([]models.Track, error)
	SetLike(trackID, userID uint32) (bool, error)
	UnLike(trackID, userID uint32) (bool, error)
	IsLiked(trackID, userID uint32) (bool, error)
}

// Repository includes DBMS-relatable methods to work with tracks
type Repository interface {
	Insert(track models.Track, artistsID []uint32) (uint32, error)
	GetByID(trackID uint32) (*models.Track, error)
	DeleteByID(trackID uint32) error
	GetFeed() ([]models.Track, error)
	GetByAlbum(albumID uint32) ([]models.Track, error)
	GetByPlaylist(playlistID uint32) ([]models.Track, error)
	GetByArtist(artistID uint32) ([]models.Track, error)
	GetLikedByUser(userID uint32) ([]models.Track, error)
	InsertLike(trackID, userID uint32) (bool, error)
	DeleteLike(trackID, userID uint32) (bool, error)
	IsLiked(trackID, userID uint32) (bool, error)
}

// Tables includes methods which return needed tables
// to work with tracks on repository-layer
type Tables interface {
	Tracks() string
	ArtistsTracks() string
	PlaylistsTracks() string
	LikedTracks() string
}
