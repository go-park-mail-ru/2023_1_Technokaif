package postgresql

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgreSQL tables
var PostgresTables = struct {
	Users         string
	Tracks        string
	Artists       string
	Albums        string
	ArtistsAlbums string
	ArtistsTracks string
	LikedAlbums   string
	LikedArtists  string
	LikedTracks   string
}{
	Users:         "Users",
	Tracks:        "Tracks",
	Artists:       "Artists",
	Albums:        "Albums",
	ArtistsAlbums: "Artists_Albums",
	ArtistsTracks: "Artists_Tracks",
	LikedAlbums:   "Liked_albums",
	LikedArtists:  "Liked_artists",
	LikedTracks:   "Liked_tracks",
}

const (
	maxIdleConns = 10
	maxOpenConns = 10
)

// Config includes info about postgres DB we want to connect to
type PostgresConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	DBSSLMode  string
}

// InitConfig inits DB configuration from environment variables
func InitPostgresConfig() (PostgresConfig, error) { // TODO CHECK FIELDS
	cfg := PostgresConfig{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if strings.TrimSpace(cfg.DBHost) == "" ||
		strings.TrimSpace(cfg.DBPort) == "" ||
		strings.TrimSpace(cfg.DBUser) == "" ||
		strings.TrimSpace(cfg.DBName) == "" ||
		strings.TrimSpace(cfg.DBPassword) == "" ||
		strings.TrimSpace(cfg.DBSSLMode) == "" {

		return PostgresConfig{}, errors.New("invalid db config")
	}

	return cfg, nil
}

// NewPostgresDB connects to chosen postgreSQL database
// and returns interaction interface of the database
func InitPostgresDB() (*sqlx.DB, error) {
	cfg, err := InitPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("can't init postgresql: %w", err)
	}

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode)

	db, err := sqlx.Open("postgres", dbInfo)
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
			return nil, fmt.Errorf("can't close postgresql (%w) after failed ping: %s", errClose, err.Error())
		}
		return nil, err
	}

	return db, nil
}
