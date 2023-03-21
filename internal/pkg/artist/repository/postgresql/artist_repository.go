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

// PostgreSQL implements artist.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{db: db, logger: l}
}

func (p *PostgreSQL) Insert(artist models.Artist) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (name, avatar_src) 
		VALUES ($1, $2);`,
		db.PostgresTables.Artists)

	if _, err := p.db.Exec(query, artist.Name, artist.AvatarSrc); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetByID(artistID uint32) (*models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT id, name, avatar_src 
		FROM %s 
		WHERE id = $1;`,
		db.PostgresTables.Artists)

	var artist models.Artist

	err := p.db.Get(&artist, query, artistID);
	if errors.Is(err, sql.ErrNoRows) {
		return &models.Artist{}, fmt.Errorf("(repo) %w: %v", &models.NoSuchArtistError{ArtistID: artistID}, err)
	} else if err != nil {
		return &models.Artist{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &artist, nil
}

func (p *PostgreSQL) Update(artist models.Artist) error {
	query := fmt.Sprintf(
		`UPDATE %s 
		SET name = $1, avatar_src = $2 
		WHERE id = $3;`,
		db.PostgresTables.Artists)

	if _, err := p.db.Exec(query, artist.Name, artist.AvatarSrc, artist.ID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) DeleteByID(artistID uint32) error {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1;`,
		db.PostgresTables.Artists)

	if _, err := p.db.Exec(query, artistID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetFeed() ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT id, name, avatar_src  
		FROM %s 
		LIMIT 100;`,
		db.PostgresTables.Artists)

	var artists []models.Artist
	if err := p.db.Select(&artists, query); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetByAlbum(albumID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src 
		FROM %s a 
			INNER JOIN %s aa ON a.id = aa.artist_id 
		WHERE aa.album_id = $1;`,
		db.PostgresTables.Artists, db.PostgresTables.ArtistsAlbums)

	var artists []models.Artist
	if err := p.db.Select(&artists, query, albumID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetByTrack(trackID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src 
		FROM %s a 
			INNER JOIN %s at ON a.id = at.artist_id 
		WHERE at.track_id = $1;`,
		db.PostgresTables.Artists, db.PostgresTables.ArtistsTracks)

	var artists []models.Artist
	if err := p.db.Select(&artists, query, trackID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetLikedByUser(userID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src, 
		FROM %s a 
			INNER JOIN %s ua ON a.id = ua.artist_id 
		WHERE ua.user_id = $1;`,
		db.PostgresTables.Artists, db.PostgresTables.LikedArtists)

	var artists []models.Artist
	if err := p.db.Select(&artists, query, userID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}
