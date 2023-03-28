package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// PostgreSQL implements track.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables track.Tables
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, t track.Tables, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
		logger: l,
	}
}

func (p *PostgreSQL) Insert(track models.Track, artistsID []uint32) (_ uint32, err error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				err = fmt.Errorf("(repo) rollback error: %w; %w", errRollback, err)
			}
		} else {
			if errCommit := tx.Commit(); errCommit != nil {
				err = fmt.Errorf("(repo) commit error: %w", errCommit)
			}
		}
	}()

	insertTrackQuery := fmt.Sprintf(
		`INSERT INTO %s (name, album_id, album_position, cover_src, record_src) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		p.tables.Tracks())

	var trackID uint32
	row := tx.QueryRow(insertTrackQuery, track.Name, track.AlbumID, track.AlbumPosition, track.CoverSrc, track.RecordSrc)
	if err := row.Scan(&trackID); err != nil {
		return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	insertTrackArtistsQuery := fmt.Sprintf(
		`INSERT INTO %s (artist_id, track_id) 
		VALUES ($1, $2);`,
		p.tables.ArtistsTracks())

	for _, artistID := range artistsID {
		if _, err := tx.Exec(insertTrackArtistsQuery, artistID, trackID); err != nil {
			return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
		}
	}

	return trackID, nil
}

func (p *PostgreSQL) GetByID(trackID uint32) (*models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, cover_src, record_src, listens
		FROM %s 
		WHERE id = $1;`,
		p.tables.Tracks())

	var track models.Track
	err := p.db.Get(&track, query, trackID)
	if errors.Is(err, sql.ErrNoRows) {
		return &models.Track{},
			fmt.Errorf("(repo) %w: %v", &models.NoSuchTrackError{TrackID: trackID}, err)
	}
	if err != nil {
		return &models.Track{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &track, nil
}

func (p *PostgreSQL) DeleteByID(trackID uint32) error {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE id = $1;`,
		p.tables.Tracks())

	if _, err := p.db.Exec(query, trackID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetFeed() ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, cover_src, record_src, listens
		FROM %s 
		LIMIT 100;`,
		p.tables.Tracks())

	var tracks []models.Track
	if err := p.db.Select(&tracks, query); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (p *PostgreSQL) GetByAlbum(albumID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, album_position, cover_src, record_src, listens
		FROM %s
		WHERE album_id = $1
		ORDER BY album_position;`,
		p.tables.Tracks())

	var tracks []models.Track
	if err := p.db.Select(&tracks, query, albumID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (p *PostgreSQL) GetByArtist(artistID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT t.id, t.name, t.album_id, t.cover_src, t.record_src, t.listens
		FROM %s t
			INNER JOIN %s at ON t.id = at.track_id 
		WHERE at.artist_id = $1;`,
		p.tables.Tracks(), p.tables.ArtistsTracks())

	var tracks []models.Track
	if err := p.db.Select(&tracks, query, artistID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (p *PostgreSQL) GetLikedByUser(userID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT t.id, name, t.album_id, t.cover_src, t.record_src, t.listens
		FROM %s t 
			INNER JOIN %s ut ON t.id = ut.track_id 
		WHERE ut.user_id = $1;`,
		p.tables.Tracks(), p.tables.LikedTracks())

	var tracks []models.Track
	if err := p.db.Select(&tracks, query, userID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

const errorLikeExists = "unique_violation"

func (p *PostgreSQL) InsertLike(trackID, userID uint32) (bool, error) {
	insertLikeQuery := fmt.Sprintf(
		`INSERT INTO %s (track_id, user_id) 
		VALUES ($1, $2)`,
		p.tables.LikedTracks())

	if _, err := p.db.Exec(insertLikeQuery, trackID, userID); err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == errorLikeExists {
				return false, nil
			}
		}

		return false, fmt.Errorf("(repo) failed to insert: %w", err)
	}

	return true, nil
}

func (p *PostgreSQL) DeleteLike(trackID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE track_id = $1 AND user_id = $2;`,
		p.tables.LikedTracks())

	resExec, err := p.db.Exec(query, trackID, userID)
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
