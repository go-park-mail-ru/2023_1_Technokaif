package playlist

import (
	"context"
	"io"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=playlist.go -destination=mocks/mock.go

type Usecase interface {
	Create(ctx context.Context, playlist models.Playlist, usersID []uint32, userID uint32) (uint32, error)
	GetByID(ctx context.Context, playlistID uint32) (*models.Playlist, error)
	UpdateInfoAndMembers(ctx context.Context, playlist models.Playlist, usersID []uint32, userID uint32) error
	UploadCover(ctx context.Context, playlistID uint32, userID uint32, file io.ReadSeeker, fileExtension string) error
	UploadCoverWrongFormatError() error
	Delete(ctx context.Context, playlistID uint32, userID uint32) error

	AddTrack(ctx context.Context, trackID, playlistID, userID uint32) error
	DeleteTrack(ctx context.Context, trackID, playlistID, userID uint32) error

	GetFeed(ctx context.Context) ([]models.Playlist, error)
	GetByUser(ctx context.Context, userID uint32) ([]models.Playlist, error)
	GetLikedByUser(ctx context.Context, userID uint32) ([]models.Playlist, error)
	SetLike(ctx context.Context, playlistID, userID uint32) (bool, error)
	UnLike(ctx context.Context, playlistID, userID uint32) (bool, error)
	IsLiked(ctx context.Context, artistID, userID uint32) (bool, error)
}

type Repository interface {
	// Check returns models.NoSuchPlaylistError if playlist-entry with given ID doesn't exist in DB
	Check(ctx context.Context, playlistID uint32) error
	Insert(ctx context.Context, playlist models.Playlist, usersID []uint32) (uint32, error)
	GetByID(ctx context.Context, playlistID uint32) (*models.Playlist, error)
	Update(ctx context.Context, playlist models.Playlist) error
	UpdateWithMembers(ctx context.Context, playlist models.Playlist, usersID []uint32) error
	DeleteByID(ctx context.Context, playlistID uint32) error

	AddTrack(ctx context.Context, trackID, playlistID uint32) error
	DeleteTrack(ctx context.Context, trackID, playlistID uint32) error

	GetFeed(ctx context.Context, amountLimit int) ([]models.Playlist, error)
	GetByUser(ctx context.Context, userID uint32) ([]models.Playlist, error)
	GetLikedByUser(ctx context.Context, userID uint32) ([]models.Playlist, error)
	InsertLike(ctx context.Context, playlistID, userID uint32) (bool, error)
	DeleteLike(ctx context.Context, playlistID, userID uint32) (bool, error)
	IsLiked(ctx context.Context, artistID, userID uint32) (bool, error)
}

// Tables includes methods which return needed tables
// to work with playlists on repository-layer
type Tables interface {
	Playlists() string
	UsersPlaylists() string
	PlaylistsTracks() string
	LikedPlaylists() string
}
