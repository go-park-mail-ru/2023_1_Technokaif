package artist

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=artist.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with artists
type Usecase interface {
	Create(ctx context.Context, artist models.Artist) (uint32, error)
	GetByID(ctx context.Context, artistID uint32) (*models.Artist, error)
	Delete(ctx context.Context, artistID uint32, userID uint32) error
	GetFeedTop(ctx context.Context, days uint32) ([]models.Artist, error)
	GetFeed(ctx context.Context) ([]models.Artist, error)
	GetByAlbum(ctx context.Context, albumID uint32) ([]models.Artist, error)
	GetByTrack(ctx context.Context, trackID uint32) ([]models.Artist, error)
	GetLikedByUser(ctx context.Context, userID uint32) ([]models.Artist, error)
	SetLike(ctx context.Context, artistID, userID uint32) (bool, error)
	UnLike(ctx context.Context, artistID, userID uint32) (bool, error)
	IsLiked(ctx context.Context, artistID, userID uint32) (bool, error)
	UpdateMonthListensPerUser(ctx context.Context) error
}

// Repository includes DBMS-relatable methods to work with artists
type Repository interface {
	// Check returns models.NoSuchArtistError if album-entry with given ID doesn't exist in DB
	Check(ctx context.Context, artistID uint32) error

	// Insert creates new entry of artist in DB with given model
	Insert(ctx context.Context, artist models.Artist) (uint32, error)

	// GetByID returns one entry of artist in DB with given ID
	GetByID(ctx context.Context, artistID uint32) (*models.Artist, error)

	// DeleteByID deletes one entry of artist with given ID
	DeleteByID(ctx context.Context, artistID uint32) error

	GetFeedTop(ctx context.Context, days, limit uint32) ([]models.Artist, error)

	// GetFeed returns artist entries with biggest amount of likes per some duration
	GetFeed(ctx context.Context, limit uint32) ([]models.Artist, error)

	// GetByAlbum returns all artist entries related with album entry with given ID
	GetByAlbum(ctx context.Context, albumID uint32) ([]models.Artist, error)

	// GetByTrack returns all artist entries related with Track with given ID
	GetByTrack(ctx context.Context, trackID uint32) ([]models.Artist, error)

	// GetByAlbum returns all Artist entries with like entry of user with given ID
	GetLikedByUser(ctx context.Context, userID uint32) ([]models.Artist, error)

	InsertLike(ctx context.Context, artistID, userID uint32) (bool, error)

	DeleteLike(ctx context.Context, artistID, userID uint32) (bool, error)

	IsLiked(ctx context.Context, artistID, userID uint32) (bool, error)

	UpdateMonthListensPerUser(ctx context.Context) error
}

// Tables includes methods which return needed tables
// to work with artists on repository layer
type Tables interface {
	Artists() string
	ArtistsAlbums() string
	ArtistsTracks() string
	LikedArtists() string
	Listens() string
	Tracks() string
}
