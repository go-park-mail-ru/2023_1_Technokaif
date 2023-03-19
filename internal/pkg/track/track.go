package track

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// TrackUsecase includes bussiness logics methods to work with albums
type TrackUsecase interface {
	GetByID(trackID uint32) (models.TrackTransfer, error)
	GetFeed() ([]models.TrackTransfer, error)
	GetByAlbum(albumID uint32) ([]models.TrackTransfer, error)
}

// TrackRepository includes DBMS-relatable methods to work with tracks
type TrackRepository interface {
	Insert(track models.Track) error
	GetByID(trackID uint32) (models.Track, error)
	Update(track models.Track) error
	Delete(id uint32) error
	GetFeed() ([]models.Track, error)
	GetByArtist(artistID uint32) ([]models.Track, error)
	GetByAlbum(albumID uint32) ([]models.Track, error)
	GetByUser(userID uint32) ([]models.Track, error)
	// GetListens(trackID uint32) (uint64, error)
	// IncrementListens(trackID uint32) error
}
