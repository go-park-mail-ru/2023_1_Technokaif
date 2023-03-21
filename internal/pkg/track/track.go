package track

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// Usecase includes bussiness logics methods to work with tracks
type Usecase interface {
	Create(track models.Track, artistsID []uint32) (uint32, error)
	GetByID(trackID uint32) (*models.Track, error)
	Change(track models.Track) error
	DeleteByID(trackID uint32) error
	GetFeed() ([]models.Track, error)
	GetByAlbum(albumID uint32) ([]models.Track, error)
	GetByArtist(artistID uint32) ([]models.Track, error)
	GetLikedByUser(userID uint32) ([]models.Track, error)
	// GetListens(trackID uint32) (uint64, error)
	// IncrementListens(trackID uint32) error
}

// Repository includes DBMS-relatable methods to work with tracks
type Repository interface {
	Insert(track models.Track, artistsID []uint32) (uint32, error)
	GetByID(trackID uint32) (*models.Track, error)
	Update(track models.Track) error
	DeleteByID(trackID uint32) error
	GetFeed() ([]models.Track, error)
	GetByAlbum(albumID uint32) ([]models.Track, error)
	GetByArtist(artistID uint32) ([]models.Track, error)
	GetLikedByUser(userID uint32) ([]models.Track, error)
	// GetListens(trackID uint32) (uint64, error)
	// IncrementListens(trackID uint32) error
}
