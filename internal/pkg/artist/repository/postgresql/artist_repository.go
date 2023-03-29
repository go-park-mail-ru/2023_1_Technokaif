package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// PostgreSQL implements artist.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables artist.Tables
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, t artist.Tables, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
		logger: l,
	}
}

func (p *PostgreSQL) Insert(artist models.Artist) (uint32, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, name, avatar_src) 
		VALUES ($1, $2, $3) RETURNING id;`,
		p.tables.Artists())

	var artistID uint32
	row := p.db.QueryRow(query, artist.UserID, artist.Name, artist.AvatarSrc)
	if err := row.Scan(&artistID); err != nil {
		return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artistID, nil
}

func (p *PostgreSQL) GetByID(artistID uint32) (*models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT id, user_id, name, avatar_src 
		FROM %s 
		WHERE id = $1;`,
		p.tables.Artists())

	var artist models.Artist

	err := p.db.Get(&artist, query, artistID)
	if errors.Is(err, sql.ErrNoRows) {
		return &models.Artist{},
			fmt.Errorf("(repo) %w: %v", &models.NoSuchArtistError{ArtistID: artistID}, err)
	}
	if err != nil {
		return &models.Artist{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &artist, nil
}

func (p *PostgreSQL) Update(artist models.Artist) error {
	query := fmt.Sprintf(
		`UPDATE %s 
		SET name = $1, avatar_src = $2 
		WHERE id = $3;`,
		p.tables.Artists())

	if _, err := p.db.Exec(query, artist.Name, artist.AvatarSrc, artist.ID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) DeleteByID(artistID uint32) error {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE id = $1;`,
		p.tables.Artists())

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
		p.tables.Artists())

	var artists []models.Artist
	if err := p.db.Select(&artists, query); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetByAlbum(albumID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.user_id, a.name, a.avatar_src 
		FROM %s a 
			INNER JOIN %s aa ON a.id = aa.artist_id 
		WHERE aa.album_id = $1;`,
		p.tables.Artists(), p.tables.ArtistsAlbums())

	var artists []models.Artist
	if err := p.db.Select(&artists, query, albumID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetByTrack(trackID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.user_id, a.name, a.avatar_src 
		FROM %s a 
			INNER JOIN %s at ON a.id = at.artist_id 
		WHERE at.track_id = $1;`,
		p.tables.Artists(), p.tables.ArtistsTracks())

	var artists []models.Artist
	if err := p.db.Select(&artists, query, trackID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetLikedByUser(userID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src
		FROM %s a 
			INNER JOIN %s ua ON a.id = ua.artist_id 
		WHERE ua.user_id = $1;`,
		p.tables.Artists(), p.tables.LikedArtists())

	var artists []models.Artist
	if err := p.db.Select(&artists, query, userID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

const errorLikeExists = "unique_violation"

func (p *PostgreSQL) InsertLike(artistID, userID uint32) (bool, error) {
	insertLikeQuery := fmt.Sprintf(
		`INSERT INTO %s (artist_id, user_id) 
		VALUES ($1, $2)`,
		p.tables.LikedArtists())

	if _, err := p.db.Exec(insertLikeQuery, artistID, userID); err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == errorLikeExists {
				return false, nil
			} 
		} 

		return false, fmt.Errorf("(repo) failed to insert: %w", err)
	}

	return true, nil
}

func (p *PostgreSQL) DeleteLike(artistID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE artist_id = $1 AND user_id = $2;`,
		p.tables.LikedArtists())

	resExec, err := p.db.Exec(query, artistID, userID)
	if err != nil {
		return false, fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("(repo) failed to check query result: %w", err)
	}

	if deleted == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
