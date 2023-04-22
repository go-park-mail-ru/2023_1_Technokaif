package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// PostgreSQL implements album.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables album.Tables
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, t album.Tables, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
		logger: l,
	}
}

func (p *PostgreSQL) Insert(album models.Album, artistsID []uint32) (_ uint32, err error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	insertAlbumQuery := fmt.Sprintf(
		`INSERT INTO %s (name, description, cover_src)
		VALUES ($1, $2, $3) RETURNING id;`,
		p.tables.Albums())

	var albumID uint32
	row := tx.QueryRow(insertAlbumQuery, album.Name, album.Description, album.CoverSrc)
	if err := row.Scan(&albumID); err != nil {
		return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	insertAlbumArtistsQuery := fmt.Sprintf(
		`INSERT INTO %s (artist_id, album_id)
		VALUES ($1, $2);`,
		p.tables.ArtistsAlbums())

	for _, artistID := range artistsID {
		if _, err := tx.Exec(insertAlbumArtistsQuery, artistID, albumID); err != nil {
			return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
		}
	}

	return albumID, nil
}

func (p *PostgreSQL) GetByID(albumID uint32) (*models.Album, error) {
	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src 
		FROM %s 
		WHERE id = $1;`,
		p.tables.Albums())

	var album models.Album

	if err := p.db.Get(&album, query, albumID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchAlbumError{AlbumID: albumID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &album, nil
}

func (p *PostgreSQL) DeleteByID(albumID uint32) error {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE id = $1;`,
		p.tables.Albums())

	resExec, err := p.db.Exec(query, albumID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return fmt.Errorf("(repo) failed to check RowsAffected: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("(repo): %w", &models.NoSuchAlbumError{AlbumID: albumID})
	}

	return nil
}

func (p *PostgreSQL) GetFeed() ([]models.Album, error) {
	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src  
		FROM %s 
		LIMIT 100;`,
		p.tables.Albums())

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
		p.tables.Albums(), p.tables.ArtistsAlbums())

	var albums []models.Album
	if err := p.db.Select(&albums, query, artistID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchArtistError{ArtistID: artistID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}

func (p *PostgreSQL) GetByTrack(trackID uint32) (*models.Album, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.description, a.cover_src 
		FROM %s a
			INNER JOIN %s t ON a.id = t.album_id
		WHERE t.id = $1;`,
		p.tables.Albums(), p.tables.Tracks())

	var album models.Album
	if err := p.db.Get(&album, query, trackID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchTrackError{TrackID: trackID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &album, nil
}

func (p *PostgreSQL) GetLikedByUser(userID uint32) ([]models.Album, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.description, a.cover_src
		FROM %s a 
			INNER JOIN %s ua ON a.id = ua.album_id 
		WHERE ua.user_id = $1;`,
		p.tables.Albums(), p.tables.LikedAlbums())

	var albums []models.Album
	if err := p.db.Select(&albums, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}

const errorLikeExists = "unique_violation"

func (p *PostgreSQL) InsertLike(albumID, userID uint32) (bool, error) {
	insertLikeQuery := fmt.Sprintf(
		`INSERT INTO %s (album_id, user_id) 
		VALUES ($1, $2)`,
		p.tables.LikedAlbums())

	if _, err := p.db.Exec(insertLikeQuery, albumID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("(repo) %w: %w", &models.NoSuchAlbumError{AlbumID: albumID}, err)
		}

		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == errorLikeExists {
				return false, nil
			}
		}

		return false, fmt.Errorf("(repo) failed to insert: %w", err)
	}

	return true, nil
}

func (p *PostgreSQL) DeleteLike(albumID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE album_id = $1 AND user_id = $2;`,
		p.tables.LikedAlbums())

	resExec, err := p.db.Exec(query, albumID, userID)
	if err != nil {
		return false, fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("(repo) failed to check query result: %w", err)
	}

	if deleted == 0 {
		return false, nil
	}
	return true, nil
}

func (p *PostgreSQL) IsLiked(albumID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`SELECT CASE WHEN 
			EXISTS(SELECT *
				FROM %s
				WHERE album_id = $1 AND user_id = $2
			) THEN TRUE ELSE FALSE END;`,
		p.tables.LikedAlbums())

	var isLiked bool
	err := p.db.Get(&isLiked, query, albumID, userID)
	if err != nil {
		return false, fmt.Errorf("(repo) failed to check if album is liked by user: %w", err)
	}

	return isLiked, nil
}
