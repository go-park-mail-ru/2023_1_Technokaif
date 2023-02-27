package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

// Usecase implements all current app's services
type Usecase struct {
	Auth
	Album
	Artist
	Track

	log *logger.Logger
}

// Auth describes which methods have to be implemented by auth-service
type Auth interface {

	// CreateUser creates new entity of user and returns it's id
	CreateUser(user models.User) (int, error)

	// GetUserID gets User's ID if such User exists
	GetUserID(username, password string) (uint, error)

	// GenerateToken returns token created with user's username and password
	GenerateAccessToken(userID uint) (string, error)

	CheckAccessToken(accessToken string) (uint, error)
}

type Album interface {
	GetAlbums() []Album
}

type Artist interface {
	GetArtists() []Artist
}

type Track interface {
	GetTracks() []Track
}

func NewUsecase(r *repository.Repository, l *logger.Logger) *Usecase {
	return &Usecase{
		Auth:   NewAuthUsecase(r.Auth),
		Album:  NewAlbumUsecase(r.Album),
		Artist: NewArtistUsecase(r.Artist),
		Track:  NewTrackUsecase(r.Track),
		log:    l,
	}
}
