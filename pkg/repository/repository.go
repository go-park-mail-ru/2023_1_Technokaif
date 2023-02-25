package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2023_1_Technokaif/models"
)

// Repository
type Repository struct {
	Auth
	// Other services
}

// Auth includes DBMS-relatable methods for authentication
type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

// NewRepository initialize SQL DBMS
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Auth: NewAuthPostgres(db),
		// Other services
	}
}
