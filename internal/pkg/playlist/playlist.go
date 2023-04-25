package playlist

import (
	"io"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=playlist.go -destination=mocks/mock.go

type Usecase interface {
	Create(playlist models.Playlist, usersID []uint32, userID uint32) (uint32, error)
	GetByID(playlistID uint32) (*models.Playlist, error)
	UpdateInfoAndMembers(playlist models.Playlist, usersID []uint32, userID uint32) error
	UploadCover(playlistID uint32, userID uint32, file io.ReadSeeker, fileExtension string) error
	Delete(playlistID uint32, userID uint32) error

	AddTrack(trackID, playlistID, userID uint32) error
	DeleteTrack(trackID, playlistID, userID uint32) error

	GetFeed() ([]models.Playlist, error)
	GetByUser(userID uint32) ([]models.Playlist, error)
	GetLikedByUser(userID uint32) ([]models.Playlist, error)
	SetLike(playlistID, userID uint32) (bool, error)
	UnLike(playlistID, userID uint32) (bool, error)
	IsLiked(artistID, userID uint32) (bool, error)
}

type Repository interface {
	// Check returns models.NoSuchPlaylistError if playlist-entry with given ID exists in DB
	Check(playlistID uint32) error
	Insert(playlist models.Playlist, usersID []uint32) (uint32, error)
	GetByID(playlistID uint32) (*models.Playlist, error)
	Update(playlist models.Playlist) error
	UpdateWithMembers(playlist models.Playlist, usersID []uint32) error
	DeleteByID(playlistID uint32) error

	AddTrack(trackID, playlistID uint32) error
	DeleteTrack(trackID, playlistID uint32) error

	GetFeed(amountLimit int) ([]models.Playlist, error)
	GetByUser(userID uint32) ([]models.Playlist, error)
	GetLikedByUser(userID uint32) ([]models.Playlist, error)
	InsertLike(playlistID, userID uint32) (bool, error)
	DeleteLike(playlistID, userID uint32) (bool, error)
	IsLiked(artistID, userID uint32) (bool, error)
}

// Tables includes methods which return needed tables
// to work with playlists on repository-layer
type Tables interface {
	Playlists() string
	UsersPlaylists() string
	PlaylistsTracks() string
	LikedPlaylists() string
}
