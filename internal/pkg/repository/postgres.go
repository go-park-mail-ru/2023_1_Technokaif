package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgreSQL tables
const (
	UsersTable         = "Users"
	TracksTable        = "Tracks"
	ArtistsTable       = "Artists"
	AlbumsTable        = "Albums"
	ArtistsAlbumsTable = "Artists_Albums"
	ArtistsTracksTable = "Artists_Tracks"
)

// Config sets info about postgres DB we want to connect to
type Config struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
	SSLMode  string
}

// NewPostgresDB connects to chosen postgreSQL database
// and returns interaction interface of the database
func NewPostrgresDB(cfg Config) (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
