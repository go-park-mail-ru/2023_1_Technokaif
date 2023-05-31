package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"

	commonSQL "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/db"
)

// PostgreSQL implements track.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables track.Tables
}

func NewPostgreSQL(db *sqlx.DB, t track.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

func (p *PostgreSQL) Check(ctx context.Context, trackID uint32) error {
	query := fmt.Sprintf(
		`SELECT EXISTS(
			SELECT id
			FROM %s
			WHERE id = $1
		);`,
		p.tables.Tracks())

	var exists bool
	err := p.db.GetContext(ctx, &exists, query, trackID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	if !exists {
		return fmt.Errorf("(repo) %w: %w", &models.NoSuchTrackError{TrackID: trackID}, err)
	}

	return nil
}

func (p *PostgreSQL) Insert(ctx context.Context, track models.Track, artistsID []uint32) (_ uint32, repoErr error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer commonSQL.CheckTransaction(tx, &repoErr)

	insertTrackQuery := fmt.Sprintf(
		`INSERT INTO %s (name, album_id, album_position, cover_src, record_src, duration) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		p.tables.Tracks())

	var trackID uint32
	row := tx.QueryRowContext(ctx, insertTrackQuery, track.Name, track.AlbumID,
		track.AlbumPosition, track.CoverSrc, track.RecordSrc, track.Duration)
	if err := row.Scan(&trackID); err != nil {
		return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	insertTrackArtistsQuery := fmt.Sprintf(
		`INSERT INTO %s (artist_id, track_id) 
		VALUES ($1, $2);`,
		p.tables.ArtistsTracks())

	for _, artistID := range artistsID {
		if _, err := tx.ExecContext(ctx, insertTrackArtistsQuery, artistID, trackID); err != nil {
			return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
		}
	}

	return trackID, nil
}

func (p *PostgreSQL) GetByID(ctx context.Context, trackID uint32) (*models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, cover_src, record_src, listens, duration
		FROM %s 
		WHERE id = $1;`,
		p.tables.Tracks())

	var track models.Track
	err := p.db.GetContext(ctx, &track, query, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Track{},
				fmt.Errorf("(repo) %w: %w", &models.NoSuchTrackError{TrackID: trackID}, err)
		}

		return &models.Track{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &track, nil
}

func (p *PostgreSQL) DeleteByID(ctx context.Context, trackID uint32) error {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE id = $1;`,
		p.tables.Tracks())

	resExec, err := p.db.ExecContext(ctx, query, trackID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return fmt.Errorf("(repo) failed to check RowsAffected: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("(repo): %w", &models.NoSuchTrackError{TrackID: trackID})
	}

	return nil
}

// func (p *PostgreSQL) GetFeed(ctx context.Context, limit uint32) ([]models.Track, error) {
// 	query := fmt.Sprintf(
// 		`SELECT id, name, album_id, cover_src, record_src, listens, duration
// 		FROM %s
// 		LIMIT $1;`,
// 		p.tables.Tracks())

// 	var tracks []models.Track
// 	if err := p.db.SelectContext(ctx, &tracks, query, limit); err != nil {
// 		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
// 	}

// 	return tracks, nil
// }

func (p *PostgreSQL) GetFeed(ctx context.Context, days, limit uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT t.id, t.name, t.album_id, t.cover_src, t.record_src, t.listens, t.duration
		FROM (
			SELECT track_id, COUNT(*) AS listens_by_time
			FROM %s
			WHERE commited_at BETWEEN (current_timestamp - $1 * interval '1 day') AND current_timestamp
			GROUP BY track_id
			ORDER BY listens_by_time DESC
		) AS tbl
			INNER JOIN %s AS t ON tbl.track_id = t.id
		LIMIT $2;`,
		p.tables.Listens(), p.tables.Tracks())

	var tracks []models.Track
	if err := p.db.SelectContext(ctx, &tracks, query, strconv.Itoa(int(days)), limit); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (p *PostgreSQL) GetByAlbum(ctx context.Context, albumID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, album_position, cover_src, record_src, listens, duration
		FROM %s
		WHERE album_id = $1
		ORDER BY album_position;`,
		p.tables.Tracks())

	var tracks []models.Track
	if err := p.db.SelectContext(ctx, &tracks, query, albumID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchAlbumError{AlbumID: albumID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (p *PostgreSQL) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT t.id, t.name, t.album_id, t.cover_src, t.record_src, t.listens, t.duration
		FROM %s t
			INNER JOIN %s pt ON t.id = pt.track_id 
		WHERE pt.playlist_id = $1
		ORDER BY pt.added_at DESC;`,
		p.tables.Tracks(), p.tables.PlaylistsTracks())

	var tracks []models.Track
	if err := p.db.SelectContext(ctx, &tracks, query, playlistID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchPlaylistError{PlaylistID: playlistID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (p *PostgreSQL) GetByArtist(ctx context.Context, artistID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT t.id, t.name, t.album_id, t.cover_src, t.record_src, t.listens, t.duration
		FROM %s t
			INNER JOIN %s at ON t.id = at.track_id 
		WHERE at.artist_id = $1;`,
		p.tables.Tracks(), p.tables.ArtistsTracks())

	var tracks []models.Track
	if err := p.db.SelectContext(ctx, &tracks, query, artistID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchArtistError{ArtistID: artistID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (p *PostgreSQL) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT t.id, name, t.album_id, t.cover_src, t.record_src, t.listens, t.duration
		FROM %s t 
			INNER JOIN %s ut ON t.id = ut.track_id 
		WHERE ut.user_id = $1
		ORDER BY liked_at DESC;`,
		p.tables.Tracks(), p.tables.LikedTracks())

	var tracks []models.Track
	if err := p.db.SelectContext(ctx, &tracks, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

const errorLikeExists = "unique_violation"

func (p *PostgreSQL) InsertLike(ctx context.Context, trackID, userID uint32) (bool, error) {
	insertLikeQuery := fmt.Sprintf(
		`INSERT INTO %s (track_id, user_id) 
		VALUES ($1, $2)`,
		p.tables.LikedTracks())

	if _, err := p.db.ExecContext(ctx, insertLikeQuery, trackID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("(repo) %w: %w", &models.NoSuchTrackError{TrackID: trackID}, err)
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

func (p *PostgreSQL) DeleteLike(ctx context.Context, trackID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE track_id = $1 AND user_id = $2;`,
		p.tables.LikedTracks())

	resExec, err := p.db.ExecContext(ctx, query, trackID, userID)
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

func (p *PostgreSQL) IsLiked(ctx context.Context, trackID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`SELECT EXISTS(
			SELECT track_id
			FROM %s
			WHERE track_id = $1 AND user_id = $2
		);`,
		p.tables.LikedTracks())

	var isLiked bool
	if err := p.db.GetContext(ctx, &isLiked, query, trackID, userID); err != nil {
		return false, fmt.Errorf("(repo) failed to check if track is liked by user: %w", err)
	}

	return isLiked, nil
}

func (p *PostgreSQL) IncrementListens(ctx context.Context, trackID, userID uint32) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (track_id, user_id)
		VALUES ($1, $2);`,
		p.tables.Listens())

	if _, err := p.db.ExecContext(ctx, query, trackID, userID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetListens(ctx context.Context, trackID uint32) (uint32, error) {
	query := fmt.Sprintf(
		`SELECT COUNT(track_id)
		FROM %s
		WHERE track_id = $1;`,
		p.tables.Listens())

	var listens uint32
	err := p.db.GetContext(ctx, &listens, query, trackID)
	if err != nil {
		return 0, fmt.Errorf("(repo) failed to get listens of track: %w", err)
	}

	return listens, nil
}

func (p *PostgreSQL) GetListensByInterval(ctx context.Context,
	start, end time.Time, trackID uint32) (uint32, error) {

	query := fmt.Sprintf(
		`SELECT COUNT(track_id)
		FROM %s
		WHERE track_id = $1
			AND commited_at BETWEEN($2, $3);`,
		p.tables.Listens())

	var listens uint32
	err := p.db.GetContext(ctx, &listens, query, trackID, start.Unix(), end.Unix())
	if err != nil {
		return 0, fmt.Errorf("(repo) failed to get listens of track: %w", err)
	}

	return listens, nil
}

func (p *PostgreSQL) UpdateListens(ctx context.Context, trackID uint32) error {
	listens, err := p.GetListens(ctx, trackID)
	if err != nil {
		return fmt.Errorf("(repo) failed to get amount of listens of track: %w", err)
	}

	query := fmt.Sprintf(
		`UPDATE %s
		SET listens = $1
		WHERE id = $2;`,
		p.tables.Listens())

	if _, err := p.db.ExecContext(ctx, query, listens, trackID); err != nil {
		return fmt.Errorf("(repo) failed to update amount of listens of track: %w", err)
	}

	return nil
}

func (p *PostgreSQL) UpdateAllListens(ctx context.Context) error {
	query := fmt.Sprintf(
		`WITH new_listens_by_track AS (
			SELECT track_id AS id, COUNT(track_id) AS new_listens
			FROM %s
			WHERE commited_at > current_timestamp - time '00:01'
			GROUP BY track_id
		)
		UPDATE %s AS t
		SET listens = listens + nlt.new_listens
		FROM new_listens_by_track AS nlt
		WHERE t.id = nlt.id;`,
		p.tables.Listens(), p.tables.Tracks())

	if _, err := p.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("(repo) failed to update all listens: %w", err)
	}

	return nil
}
