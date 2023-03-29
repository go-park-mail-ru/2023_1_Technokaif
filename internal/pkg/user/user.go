package user

import (
	"io"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=user.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with users
type Usecase interface {
	GetByID(userID uint32) (*models.User, error)
	UpdateInfo(user *models.User) error
	UploadAvatar(user *models.User, file io.ReadSeeker, fileExtension string) error
	UploadAvatarWrongFormatError() error
}

// Repository includes DBMS-relatable methods to work with users
type Repository interface {
	GetByID(userID uint32) (*models.User, error)

	// CreateUser inserts new user into DB and return it's id
	// or error if it already exists
	CreateUser(user models.User) (uint32, error)

	UpdateInfo(user *models.User) error

	// GetUserByUsername returns models.User if it's entry in DB exists or error otherwise
	GetUserByUsername(username string) (*models.User, error)

	// GetFriends returns all users, who have friendship entry with user with given ID
	// GetFriends(userID uint32) ([]models.User, error)

	// GetListenedAlbum(albumID uint32) ([]models.User, error)
	// GetListenedTrack(trackID uint32) ([]models.User, error)
	// GetLikedAlbum(albumID uint32) ([]models.User, error)
	// GetLikedArtist(artistID uint32) ([]models.User, error)
	// GetLikedTrack(trackID uint32) ([]models.User, error)
}

// Tables includes methods which return needed tables
// to work with users on repository-layer
type Tables interface {
	Users() string
}
