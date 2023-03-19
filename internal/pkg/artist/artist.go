package artist

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// ArtistUsecase includes bussiness logics methods to work with albums
type ArtistUsecase interface {
	GetByID(artistID uint32) (models.ArtistTransfer, error)
	GetFeed() ([]models.ArtistTransfer, error)
}

// ArtistRepository includes DBMS-relatable methods to work with artists
type ArtistRepository interface {
	GetByID(artistID uint32) (models.Artist, error)
	GetFeed() ([]models.Artist, error)
	GetByTrack(trackID uint32) ([]models.Artist, error)
	GetByAlbum(albumID uint32) ([]models.Artist, error)
	// GetListens(artistID uint32) (uint32, error)
}
