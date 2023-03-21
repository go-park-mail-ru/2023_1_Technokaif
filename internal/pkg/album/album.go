package album

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// Usecase includes bussiness logics methods to work with albums
type Usecase interface {
	Create(album models.Album) error
	GetByID(albumID uint32) (*models.Album, error)
	Change(album models.Album) error
	DeleteByID(albumID uint32) error
	GetFeed() ([]models.Album, error)
	GetByArtist(artistID uint32) ([]models.Album, error)
	GetByTrack(trackID uint32) (*models.Album, error)
	GetLikedByUser(userID uint32) ([]models.Album, error)
	// GetListens(albumID uint32) (uint32, error)
}

// Repository includes DBMS-relatable methods to work with albums
type Repository interface {
	Insert(album models.Album) error
	GetByID(albumID uint32) (*models.Album, error)
	Update(album models.Album) error
	DeleteByID(albumID uint32) error
	GetFeed() ([]models.Album, error)
	GetByArtist(artistID uint32) ([]models.Album, error)
	GetByTrack(trackID uint32) (*models.Album, error)
	GetLikedByUser(userID uint32) ([]models.Album, error)
	// GetListens(albumID uint32) (uint32, error)
}
