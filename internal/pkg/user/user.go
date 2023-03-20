package user

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// Usecase includes bussiness logics methods to work with users
type Usecase interface {
	GetByID(userID uint32) (models.UserTransfer, error)
}

// Repository includes DBMS-relatable methods to work with users
type Repository interface {
	GetByID(userID uint32) (models.User, error)

	// GetFriends returns all users, who have friendship entry with user with given ID
	// GetFriends(userID uint32) ([]models.User, error)

	// GetListenedAlbum(albumID uint32) ([]models.User, error)
	// GetListenedTrack(trackID uint32) ([]models.User, error)
	// GetLikedAlbum(albumID uint32) ([]models.User, error)
	// GetLikedArtist(artistID uint32) ([]models.User, error)
	// GetLikedTrack(trackID uint32) ([]models.User, error)
}
