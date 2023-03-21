package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	db "github.com/go-park-mail-ru/2023_1_Technokaif/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// PostgreSQL implements album.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{db: db, logger: l}
}

func (p *PostgreSQL) Insert(album models.Album) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (name, description, cover_src)
		VALUES ($1, $2, $3);`,
		db.PostgresTables.Albums)

	if _, err := p.db.Exec(query, album.Name, album.Description, album.CoverSrc); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetByID(albumID uint32) (*models.Album, error) {
	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src 
		FROM %s 
		WHERE id = $1;`,
		db.PostgresTables.Albums)

	var albums models.Album

	err := p.db.Get(&albums, query, albumID); 
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("(repo) %w: %v", &models.NoSuchAlbumError{AlbumID: albumID}, err)
	} else if err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &albums, nil
}

func (p *PostgreSQL) Update(album models.Album) error {
	query := fmt.Sprintf(
		`UPDATE %s 
		SET name = $1, description = $2, cover_src = $3 
		WHERE id = $4;`,
		db.PostgresTables.Albums)

	if _, err := p.db.Exec(query, album.Name, album.Description,
		album.CoverSrc, album.ID); err != nil {

		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) DeleteByID(albumID uint32) error {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1;`,
		db.PostgresTables.Albums)

	if _, err := p.db.Exec(query, albumID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetFeed() ([]models.Album, error) {
	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src  
		FROM %s 
		LIMIT 100;`,
		db.PostgresTables.Albums)

	var albums []models.Album
	if err := p.db.Select(&albums, query); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}

func (p *PostgreSQL) GetByArtist(artistID uint32) ([]models.Album, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.description, a.cover_src 
		FROM %s a
			INNER JOIN %s aa ON a.id = aa.album_id
		WHERE aa.artist_id = $1;`,
		db.PostgresTables.Albums, db.PostgresTables.ArtistsAlbums)

	var albums []models.Album
	if err := p.db.Select(&albums, query, artistID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}

func (p *PostgreSQL) GetByTrack(trackID uint32) (*models.Album, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.description, a.cover_src 
		FROM %s a
			INNER JOIN %s t ON a.id = t.album_id
		WHERE t.album_id = $1;`,
		db.PostgresTables.Albums, db.PostgresTables.Tracks)

	var album models.Album
	if err := p.db.Select(&album, query, trackID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &album, nil
}

func (p *PostgreSQL) GetLikedByUser(userID uint32) ([]models.Album, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.description, a.cover_src, 
		FROM %s a 
			INNER JOIN %s ua ON a.id = ua.artist_id 
		WHERE ua.user_id = $1;`,
		db.PostgresTables.Albums, db.PostgresTables.LikedAlbums)

	var albums []models.Album
	if err := p.db.Select(&albums, query, userID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}
