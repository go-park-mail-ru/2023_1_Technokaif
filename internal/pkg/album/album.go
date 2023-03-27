package album

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// Usecase includes bussiness logics methods to work with albums
type Usecase interface {
	Create(album models.Album, artistsID []uint32, userID uint32) (uint32, error)
	GetByID(albumID uint32) (*models.Album, error)
	Change(album models.Album) error
	Delete(albumID uint32, userID uint32) error
	GetFeed() ([]models.Album, error)
	GetByArtist(artistID uint32) ([]models.Album, error)
	GetByTrack(trackID uint32) (*models.Album, error)
	GetLikedByUser(userID uint32) ([]models.Album, error)
	SetLike(albumID, userID uint32) (bool, error)
	UnLike(albumID, userID uint32) (bool, error)
	// GetListens(albumID uint32) (uint32, error)
}

// Repository includes DBMS-relatable methods to work with albums
type Repository interface {
	Insert(album models.Album, artistsID []uint32) (uint32, error)
	GetByID(albumID uint32) (*models.Album, error)
	Update(album models.Album) error
	DeleteByID(albumID uint32) error
	GetFeed() ([]models.Album, error)
	GetByArtist(artistID uint32) ([]models.Album, error)
	GetByTrack(trackID uint32) (*models.Album, error)
	GetLikedByUser(userID uint32) ([]models.Album, error)
	InsertLike(albumID, userID uint32) (bool, error)
	DeleteLike(albumID, userID uint32) (bool, error)
	// GetListens(albumID uint32) (uint32, error)
}

// Tables includes methods which return needed tables
// to work with albums on repository-layer
type Tables interface {
	Albums() string
	Tracks() string
	ArtistsAlbums() string
	LikedAlbums() string
}
