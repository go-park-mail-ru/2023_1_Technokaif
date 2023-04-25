package album

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

//go:generate mockgen -source=album.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with albums
type Usecase interface {
	Create(album models.Album, artistsID []uint32, userID uint32) (uint32, error)
	GetByID(albumID uint32) (*models.Album, error)
	Delete(albumID uint32, userID uint32) error
	GetFeed() ([]models.Album, error)
	GetByArtist(artistID uint32) ([]models.Album, error)
	GetByTrack(trackID uint32) (*models.Album, error)
	GetLikedByUser(userID uint32) ([]models.Album, error)
	SetLike(albumID, userID uint32) (bool, error)
	UnLike(albumID, userID uint32) (bool, error)
	IsLiked(albumID, userID uint32) (bool, error)
}

// Repository includes DBMS-relatable methods to work with albums
type Repository interface {
	// Check returns models.NoSuchAlbumError if album-entry with given ID exists in DB
	Check(albumID uint32) error
	Insert(album models.Album, artistsID []uint32) (uint32, error)
	GetByID(albumID uint32) (*models.Album, error)
	DeleteByID(albumID uint32) error
	GetFeed(amountLimit int) ([]models.Album, error)
	GetByArtist(artistID uint32) ([]models.Album, error)
	GetByTrack(trackID uint32) (*models.Album, error)
	GetLikedByUser(userID uint32) ([]models.Album, error)
	InsertLike(albumID, userID uint32) (bool, error)
	DeleteLike(albumID, userID uint32) (bool, error)
	IsLiked(albumID, userID uint32) (bool, error)
}

// Tables includes methods which return needed tables
// to work with albums on repository-layer
type Tables interface {
	Albums() string
	Tracks() string
	ArtistsAlbums() string
	LikedAlbums() string
}
