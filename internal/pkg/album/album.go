package album

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// AlbumUsecase includes bussiness logics methods to work with albums
type AlbumUsecase interface {
	GetByID(albumID uint32) ([]models.AlbumTransfer, error)
	GetFeed() ([]models.AlbumTransfer, error)
}

// AlbumRepository includes DBMS-relatable methods to work with albums
type AlbumRepository interface {
	GetByID(albumID uint32) ([]models.Album, error)
	GetFeed() ([]models.Album, error)
	GetByArtist(artistID uint32) ([]models.Album, error)
	GetByTrack(trackID uint32) (models.Album, error)
	GetByUser(userID uint32) ([]models.Album, error)
	// GetListens(albumID uint32) (uint32, error)
}
