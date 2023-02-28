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
	// CreateUser inserts new user into DB and return it's id
	// or error if it already exists
	CreateUser(user models.User) (int, error)

	// GetUser returns models.User if it's entry in DB exists or error otherwise
	GetUser(username, password string) (models.User, error)
}

// Artist includes DBMS-relatable methods to work with artists
type Artist interface {
	GetFeed() ([]models.ArtistFeed, error)
}

// Album includes DBMS-relatable methods to work with albums
type Album interface {
	GetFeed() ([]models.AlbumFeed, error)
}

// Track includes DBMS-relatable methods to work with tracks
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
type UserAlreadyExistsError struct{}
type NoSuchUserError struct{}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}

func (e *NoSuchUserError) Error() string {
	return "no such user"
}
