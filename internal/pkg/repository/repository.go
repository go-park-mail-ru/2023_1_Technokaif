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

	log *logger.Logger
}

// Auth includes DBMS-relatable methods for authentication
type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Artist interface {
}

type Album interface {
}

type Track interface {
	GetFeed() ([]models.TrackFeed, error)
}

// NewRepository initialize SQL DBMS
func NewRepository(db *sql.DB, l *logger.Logger) *Repository {
	return &Repository{
		Auth:   NewAuthPostgres(db),
		Artist: NewArtistPostgres(db),
		Album:  NewAlbumPostgres(db),
		Track:  NewTrackPostgres(db),
		log:    l,
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
