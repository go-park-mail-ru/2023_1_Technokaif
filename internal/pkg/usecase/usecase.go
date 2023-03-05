package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mock.go

// Usecase implements all current app's services
type Usecase struct {
	Auth
	Album
	Artist
	Track
}

// Auth describes which methods have to be implemented by auth-service
type Auth interface {
	// LoginUser finds user in repository and returns its access token
	LoginUser(username, password string) (string, error)

	// CreateUser creates new user and returns it's id
	CreateUser(user models.User) (int, error)

	// GetUserID returns User's ID if such User exists
	GetUserByCreds(username, password string) (*models.User, error)

	// GetUserByAuthData ...
	GetUserByAuthData(userID, userVersion uint) (*models.User, error)

	// GenerateAccessToken returns token created with user's username and password
	GenerateAccessToken(userID, userVersion uint) (string, error)

	// CheckAccessToken ...
	CheckAccessToken(accessToken string) (uint, uint, error)

	// ChangeUserVersion ...
	ChangeUserVersion(userID uint) error
}

type Album interface {
	GetFeed() ([]models.AlbumFeed, error)
}

type Artist interface {
	GetFeed() ([]models.ArtistFeed, error)
}

type Track interface {
	GetFeed() ([]models.TrackFeed, error)
}

func NewUsecase(r *repository.Repository, l logger.Logger) *Usecase {
	return &Usecase{
		Auth:   NewAuthUsecase(r.Auth, l),
		Album:  NewAlbumUsecase(r.Album, l),
		Artist: NewArtistUsecase(r.Artist, l),
		Track:  NewTrackUsecase(r.Track, l),
	}
}
