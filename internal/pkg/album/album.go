package album

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=album.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with albums
type Usecase interface {
	Create(ctx context.Context, album models.Album, artistsID []uint32, userID uint32) (uint32, error)
	GetByID(ctx context.Context, albumID uint32) (*models.Album, error)
	Delete(ctx context.Context, albumID uint32, userID uint32) error
	GetFeedTop(ctx context.Context, days uint32) ([]models.Album, error)
	GetFeed(ctx context.Context) ([]models.Album, error)
	GetByArtist(ctx context.Context, artistID uint32) ([]models.Album, error)
	GetByTrack(ctx context.Context, trackID uint32) (*models.Album, error)
	GetLikedByUser(ctx context.Context, userID uint32) ([]models.Album, error)
	SetLike(ctx context.Context, albumID, userID uint32) (bool, error)
	UnLike(ctx context.Context, albumID, userID uint32) (bool, error)
	IsLiked(ctx context.Context, albumID, userID uint32) (bool, error)
}

// Repository includes DBMS-relatable methods to work with albums
type Repository interface {
	// Check returns models.NoSuchAlbumError if album-entry with given ID doesn't exist in DB
	Check(ctx context.Context, albumID uint32) error
	Insert(ctx context.Context, album models.Album, artistsID []uint32) (uint32, error)
	GetByID(ctx context.Context, albumID uint32) (*models.Album, error)
	DeleteByID(ctx context.Context, albumID uint32) error
	GetFeedTop(ctx context.Context, days, limit uint32) ([]models.Album, error)
	GetFeed(ctx context.Context, limit uint32) ([]models.Album, error)
	GetByArtist(ctx context.Context, artistID uint32) ([]models.Album, error)
	GetByTrack(ctx context.Context, trackID uint32) (*models.Album, error)
	GetLikedByUser(ctx context.Context, userID uint32) ([]models.Album, error)
	InsertLike(ctx context.Context, albumID, userID uint32) (bool, error)
	DeleteLike(ctx context.Context, albumID, userID uint32) (bool, error)
	IsLiked(ctx context.Context, albumID, userID uint32) (bool, error)
}

// Tables includes methods which return needed tables
// to work with albums on repository layer
type Tables interface {
	Albums() string
	Tracks() string
	ArtistsAlbums() string
	LikedAlbums() string
	Listens() string
}
