package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"

	commonSQL "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/db"
)

// PostgreSQL implements album.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables album.Tables
}

func NewPostgreSQL(db *sqlx.DB, t album.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

func (p *PostgreSQL) Check(ctx context.Context, albumID uint32) error {
	query := fmt.Sprintf(
		`SELECT EXISTS(
			SELECT id
			FROM %s
			WHERE id = $1
		);`,
		p.tables.Albums())

	var exists bool
	err := p.db.GetContext(ctx, &exists, query, albumID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	if !exists {
		return fmt.Errorf("(repo) %w: %w", &models.NoSuchAlbumError{AlbumID: albumID}, err)
	}

	return nil
}

func (p *PostgreSQL) Insert(ctx context.Context, album models.Album, artistsID []uint32) (_ uint32, repoErr error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer commonSQL.CheckTransaction(tx, &repoErr)

	insertAlbumQuery := fmt.Sprintf(
		`INSERT INTO %s (name, description, cover_src)
		VALUES ($1, $2, $3) RETURNING id;`,
		p.tables.Albums())

	var albumID uint32
	row := tx.QueryRowContext(ctx, insertAlbumQuery, album.Name, album.Description, album.CoverSrc)
	if err := row.Scan(&albumID); err != nil {
		return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	insertAlbumArtistsQuery := fmt.Sprintf(
		`INSERT INTO %s (artist_id, album_id)
		VALUES ($1, $2);`,
		p.tables.ArtistsAlbums())

	for _, artistID := range artistsID {
		if _, err := tx.ExecContext(ctx, insertAlbumArtistsQuery, artistID, albumID); err != nil {
			return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
		}
	}

	return albumID, nil
}

func (p *PostgreSQL) GetByID(ctx context.Context, albumID uint32) (*models.Album, error) {
	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src 
		FROM %s 
		WHERE id = $1;`,
		p.tables.Albums())

	var album models.Album

	if err := p.db.GetContext(ctx, &album, query, albumID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchAlbumError{AlbumID: albumID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &album, nil
}

func (p *PostgreSQL) DeleteByID(ctx context.Context, albumID uint32) error {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE id = $1;`,
		p.tables.Albums())

	resExec, err := p.db.ExecContext(ctx, query, albumID)
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

func (p *PostgreSQL) GetFeed(ctx context.Context, limit uint32) ([]models.Album, error) {
	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src  
		FROM %s 
		LIMIT $1;`,
		p.tables.Albums())

	var albums []models.Album
	if err := p.db.SelectContext(ctx, &albums, query, limit); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}

func (p *PostgreSQL) GetByArtist(ctx context.Context, artistID uint32) ([]models.Album, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.description, a.cover_src 
		FROM %s a
			INNER JOIN %s aa ON a.id = aa.album_id
		WHERE aa.artist_id = $1;`,
		p.tables.Albums(), p.tables.ArtistsAlbums())

	var albums []models.Album
	if err := p.db.SelectContext(ctx, &albums, query, artistID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchArtistError{ArtistID: artistID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}

func (p *PostgreSQL) GetByTrack(ctx context.Context, trackID uint32) (*models.Album, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.description, a.cover_src 
		FROM %s a
			INNER JOIN %s t ON a.id = t.album_id
		WHERE t.id = $1;`,
		p.tables.Albums(), p.tables.Tracks())

	var album models.Album
	if err := p.db.GetContext(ctx, &album, query, trackID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchTrackError{TrackID: trackID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &album, nil
}

func (p *PostgreSQL) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Album, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.description, a.cover_src
		FROM %s a 
			INNER JOIN %s ua ON a.id = ua.album_id 
		WHERE ua.user_id = $1
		ORDER BY liked_at DESC;`,
		p.tables.Albums(), p.tables.LikedAlbums())

	var albums []models.Album
	if err := p.db.SelectContext(ctx, &albums, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}

const errorLikeExists = "unique_violation"

func (p *PostgreSQL) InsertLike(ctx context.Context, albumID, userID uint32) (bool, error) {
	insertLikeQuery := fmt.Sprintf(
		`INSERT INTO %s (album_id, user_id) 
		VALUES ($1, $2)`,
		p.tables.LikedAlbums())

	if _, err := p.db.ExecContext(ctx, insertLikeQuery, albumID, userID); err != nil {
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

func (p *PostgreSQL) DeleteLike(ctx context.Context, albumID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE album_id = $1 AND user_id = $2;`,
		p.tables.LikedAlbums())

	resExec, err := p.db.ExecContext(ctx, query, albumID, userID)
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

func (p *PostgreSQL) IsLiked(ctx context.Context, albumID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`SELECT EXISTS(
			SELECT album_id
			FROM %s
			WHERE album_id = $1 AND user_id = $2
		);`,
		p.tables.LikedAlbums())

	var isLiked bool
	err := p.db.GetContext(ctx, &isLiked, query, albumID, userID)
	if err != nil {
		return false, fmt.Errorf("(repo) failed to check if album is liked by user: %w", err)
	}

	return isLiked, nil
}
