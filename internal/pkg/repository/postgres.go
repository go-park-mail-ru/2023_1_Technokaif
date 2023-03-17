package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgreSQL tables
const (
	maxIdleConns = 10
	maxOpenConns = 10

	usersTable         = "Users"
	tracksTable        = "Tracks"
	artistsTable       = "Artists"
	albumsTable        = "Albums"
	artistsAlbumsTable = "Artists_Albums"
	artistsTracksTable = "Artists_Tracks"
)

// Config includes info about postgres DB we want to connect to
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	DBSSLMode  string
}

// NewPostgresDB connects to chosen postgreSQL database
// and returns interaction interface of the database
func NewPostrgresDB(cfg Config) (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	err = db.Ping()
	if err != nil {
		errClose := db.Close()
		if errClose != nil {
			// TODO: change %w and %s (double wrap only in go1.20+)
			return nil, fmt.Errorf("can't close DB (%w) after failed ping: %s", errClose, err.Error())
		}
		return nil, err
	}

	return db, nil
}
