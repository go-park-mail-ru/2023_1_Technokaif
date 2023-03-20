package user

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

type UserUsecase interface {
	GetByID(userID uint32) (models.UserTransfer, error)
}

type UserRepository interface {
	GetByID(userID uint32) (models.User, error)

	// GetFriends returns all users, who have friendship entry with user with given ID
	// GetFriends(userID uint32) ([]models.User, error)

	// GetListenedAlbum(albumID uint32) ([]models.User, error)
	// GetListenedTrack(trackID uint32) ([]models.User, error)
	// GetLikedAlbum(albumID uint32) ([]models.User, error)
	// GetLikedArtist(artistID uint32) ([]models.User, error)
	// GetLikedTrack(trackID uint32) ([]models.User, error)
}
