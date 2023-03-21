package postgresql

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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
func initPostgresConfig() (PostgresConfig, error) { // TODO CHECK FIELDS
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
func InitPostgresDB() (*sqlx.DB, *PostgreSQLTables, error) {
	cfg, err := initPostgresConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("can't init postgresql: %w", err)
	}

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode)

	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		return nil, nil, err
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	err = db.Ping()
	if err != nil {
		errClose := db.Close()
		if errClose != nil {
			// TODO: change %w and %s (double wrap only in go1.20+)
			return nil, nil, fmt.Errorf("can't close postgresql (%w) after failed ping: %s", errClose, err.Error())
		}
		return nil, nil, err
	}

	return db, &PostgreSQLTables{}, nil
}

type PostgreSQLTables struct{}

func (pt *PostgreSQLTables) Users() string {
	return "Users"
}

func (pt *PostgreSQLTables) Tracks() string {
	return "Tracks"
}

func (pt *PostgreSQLTables) Artists() string {
	return "Artists"
}

func (pt *PostgreSQLTables) Albums() string {
	return "Albums"
}

func (pt *PostgreSQLTables) Listens() string {
	return "Listens"
}

func (pt *PostgreSQLTables) ArtistsAlbums() string {
	return "Artists_Albums"
}

func (pt *PostgreSQLTables) ArtistsTracks() string {
	return "Artists_Tracks"
}

func (pt *PostgreSQLTables) LikedAlbums() string {
	return "Liked_albums"
}

func (pt *PostgreSQLTables) LikedArtists() string {
	return "Liked_artists"
}

func (pt *PostgreSQLTables) LikedTracks() string {
	return "Liked_tracks"
}
