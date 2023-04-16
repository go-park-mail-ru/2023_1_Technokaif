package playlist

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

//go:generate mockgen -source=playlist.go -destination=mocks/mock.go

type Usecase interface {
	Create(playlist models.Playlist, usersID []uint32, userID uint32) (uint32, error)
	GetByID(playlistID uint32) (*models.Playlist, error)
	Update(playlist models.Playlist, usersID []uint32, userID uint32) error
	Delete(playlistID uint32, userID uint32) error
	GetFeed() ([]models.Playlist, error)
	GetByUser(userID uint32) ([]models.Playlist, error)
	GetLikedByUser(userID uint32) ([]models.Playlist, error)
	SetLike(playlistID, userID uint32) (bool, error)
	UnLike(playlistID, userID uint32) (bool, error)
}

type Repository interface {
	Insert(playlist models.Playlist, usersID []uint32) (uint32, error)
	GetByID(playlistID uint32) (*models.Playlist, error)
	Update(playlist models.Playlist, usersID []uint32) error
	DeleteByID(playlistID uint32) error
	GetFeed() ([]models.Playlist, error)
	GetByUser(userID uint32) ([]models.Playlist, error)
	GetLikedByUser(userID uint32) ([]models.Playlist, error)
	InsertLike(playlistID, userID uint32) (bool, error)
	DeleteLike(playlistID, userID uint32) (bool, error)
}

// Tables includes methods which return needed tables
// to work with playlists on repository-layer
type Tables interface {
	Playlists() string
	Tracks() string
	UsersPlaylists() string
	PlaylistsTracks() string
	LikedPlaylists() string
}
