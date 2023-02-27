package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Repository
type Repository struct {
	Auth
	Artist
	Album
	Track

	logger logger.Logger
}

// Auth includes DBMS-relatable methods for authentication
type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Artist interface {
	GetFeed() ([]models.ArtistFeed, error)
}

type Album interface {
	GetFeed() ([]models.AlbumFeed, error)
}

type Track interface {
	GetFeed() ([]models.TrackFeed, error)
}

// NewRepository initialize SQL DBMS
func NewRepository(db *sql.DB, l logger.Logger) *Repository {
	return &Repository{
		Auth:   NewAuthPostgres(db),
		Artist: NewArtistPostgres(db),
		Album:  NewAlbumPostgres(db),
		Track:  NewTrackPostgres(db),
		logger: l,
	}
}

// AUTH ERRORS
type UserAlreadyExistsError struct {
}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}

type NoSuchUserError struct {
}

func (e *NoSuchUserError) Error() string {
	return "no such user"
}
