package track

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=track.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with tracks
type Usecase interface {
	Create(ctx context.Context, track models.Track, artistsID []uint32, userID uint32) (uint32, error)
	GetByID(ctx context.Context, trackID uint32) (*models.Track, error)
	Delete(ctx context.Context, trackID uint32, userID uint32) error
	GetFeed(ctx context.Context) ([]models.Track, error)
	GetByAlbum(ctx context.Context, albumID uint32) ([]models.Track, error)
	GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.Track, error)
	GetByArtist(ctx context.Context, artistID uint32) ([]models.Track, error)
	GetLikedByUser(ctx context.Context, userID uint32) ([]models.Track, error)
	SetLike(ctx context.Context, trackID, userID uint32) (bool, error)
	UnLike(ctx context.Context, trackID, userID uint32) (bool, error)
	IsLiked(ctx context.Context, trackID, userID uint32) (bool, error)

	IncrementListens(ctx context.Context, trackID, userID uint32) error
	GetListens(ctx context.Context, trackID uint32) (uint32, error)
	UpdateListens(ctx context.Context, trackID uint32) error
	UpdateAllListens(ctx context.Context) error
}

// Repository includes DBMS-relatable methods to work with tracks
type Repository interface {
	// Check returns models.NoSuchTrackError if track-entry with given ID doesn't exist in DB
	Check(ctx context.Context, trackID uint32) error
	Insert(ctx context.Context, track models.Track, artistsID []uint32) (uint32, error)
	GetByID(ctx context.Context, trackID uint32) (*models.Track, error)
	DeleteByID(ctx context.Context, trackID uint32) error
	GetFeed(ctx context.Context, limit uint32) ([]models.Track, error)
	GetByAlbum(ctx context.Context, albumID uint32) ([]models.Track, error)
	GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.Track, error)
	GetByArtist(ctx context.Context, artistID uint32) ([]models.Track, error)
	GetLikedByUser(ctx context.Context, userID uint32) ([]models.Track, error)
	InsertLike(ctx context.Context, trackID, userID uint32) (bool, error)
	DeleteLike(ctx context.Context, trackID, userID uint32) (bool, error)
	IsLiked(ctx context.Context, trackID, userID uint32) (bool, error)

	IncrementListens(ctx context.Context, trackID, userID uint32) error
	GetListens(ctx context.Context, trackID uint32) (uint32, error)
	GetListensByInterval(ctx context.Context, start, end time.Time, trackID uint32) (uint32, error)
	UpdateListens(ctx context.Context, trackID uint32) error
	UpdateAllListens(ctx context.Context) error
}

// Tables includes methods which return needed tables
// to work with tracks on repository layer
type Tables interface {
	Tracks() string
	ArtistsTracks() string
	PlaylistsTracks() string
	LikedTracks() string
	Listens() string
}
