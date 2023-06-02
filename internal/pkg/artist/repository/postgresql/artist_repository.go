package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
)

// PostgreSQL implements artist.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables artist.Tables
}

func NewPostgreSQL(db *sqlx.DB, t artist.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

func (p *PostgreSQL) Check(ctx context.Context, artistID uint32) error {
	query := fmt.Sprintf(
		`SELECT EXISTS(
			SELECT id
			FROM %s
			WHERE id = $1
		);`,
		p.tables.Artists())

	var exists bool
	err := p.db.GetContext(ctx, &exists, query, artistID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	if !exists {
		return fmt.Errorf("(repo) %w: %w", &models.NoSuchArtistError{ArtistID: artistID}, err)
	}

	return nil
}

func (p *PostgreSQL) Insert(ctx context.Context, artist models.Artist) (uint32, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, name, avatar_src) 
		VALUES ($1, $2, $3) RETURNING id;`,
		p.tables.Artists())

	var artistID uint32
	row := p.db.QueryRowContext(ctx, query, artist.UserID, artist.Name, artist.AvatarSrc)
	if err := row.Scan(&artistID); err != nil {
		return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artistID, nil
}

func (p *PostgreSQL) GetByID(ctx context.Context, artistID uint32) (*models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT id, user_id, name, avatar_src, listens
		FROM %s 
		WHERE id = $1;`,
		p.tables.Artists())

	var artist models.Artist

	err := p.db.GetContext(ctx, &artist, query, artistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Artist{},
				fmt.Errorf("(repo) %w: %w", &models.NoSuchArtistError{ArtistID: artistID}, err)
		}
		return &models.Artist{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &artist, nil
}

func (p *PostgreSQL) DeleteByID(ctx context.Context, artistID uint32) error {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE id = $1;`,
		p.tables.Artists())

	resExec, err := p.db.ExecContext(ctx, query, artistID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return fmt.Errorf("(repo) failed to check affected rows: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("(repo): %w", &models.NoSuchArtistError{ArtistID: artistID})
	}

	return nil
}

func (p *PostgreSQL) GetFeedTop(ctx context.Context, days, limit uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src, a.listens
		FROM (
			SELECT track_id, COUNT(*) AS listens_by_time
			FROM %s
			WHERE commited_at BETWEEN (current_timestamp - $1 * interval '1 day') AND current_timestamp
			GROUP BY track_id
		) AS tbl
			RIGHT JOIN %s AS t ON tbl.track_id = t.id
			INNER JOIN %s AS at ON t.id = at.track_id
			INNER JOIN %s AS a on at.artist_id = a.id
		GROUP BY a.id
		ORDER BY SUM(tbl.listens_by_time) DESC NULLS LAST
		LIMIT $2;`,
		p.tables.Listens(), p.tables.Tracks(), p.tables.ArtistsTracks(), p.tables.Artists())

	var artists []models.Artist
	if err := p.db.SelectContext(ctx, &artists, query, strconv.Itoa(int(days)), limit); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetFeed(ctx context.Context, limit uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT id, name, avatar_src, listens  
		FROM %s 
		LIMIT $1;`,
		p.tables.Artists())

	var artists []models.Artist
	if err := p.db.SelectContext(ctx, &artists, query, limit); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetByAlbum(ctx context.Context, albumID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.user_id, a.name, a.avatar_src, a.listens
		FROM %s a 
			INNER JOIN %s aa ON a.id = aa.artist_id 
		WHERE aa.album_id = $1;`,
		p.tables.Artists(), p.tables.ArtistsAlbums())

	var artists []models.Artist
	if err := p.db.SelectContext(ctx, &artists, query, albumID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchAlbumError{AlbumID: albumID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetByTrack(ctx context.Context, trackID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.user_id, a.name, a.avatar_src, a.listens 
		FROM %s a 
			INNER JOIN %s at ON a.id = at.artist_id 
		WHERE at.track_id = $1;`,
		p.tables.Artists(), p.tables.ArtistsTracks())

	var artists []models.Artist
	if err := p.db.SelectContext(ctx, &artists, query, trackID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchTrackError{TrackID: trackID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src, a.listens
		FROM %s a 
			INNER JOIN %s ua ON a.id = ua.artist_id 
		WHERE ua.user_id = $1
		ORDER BY liked_at DESC;`,
		p.tables.Artists(), p.tables.LikedArtists())

	var artists []models.Artist
	if err := p.db.SelectContext(ctx, &artists, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

const errorLikeExists = "unique_violation"

func (p *PostgreSQL) InsertLike(ctx context.Context, artistID, userID uint32) (bool, error) {
	insertLikeQuery := fmt.Sprintf(
		`INSERT INTO %s (artist_id, user_id) 
		VALUES ($1, $2)`,
		p.tables.LikedArtists())

	if _, err := p.db.ExecContext(ctx, insertLikeQuery, artistID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("(repo) %w: %w", &models.NoSuchArtistError{ArtistID: artistID}, err)
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

func (p *PostgreSQL) DeleteLike(ctx context.Context, artistID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE artist_id = $1 AND user_id = $2;`,
		p.tables.LikedArtists())

	resExec, err := p.db.ExecContext(ctx, query, artistID, userID)
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

func (p *PostgreSQL) IsLiked(ctx context.Context, artistID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`SELECT EXISTS(
			SELECT artist_id
			FROM %s
			WHERE artist_id = $1 AND user_id = $2
		);`,
		p.tables.LikedArtists())

	var isLiked bool
	err := p.db.GetContext(ctx, &isLiked, query, artistID, userID)
	if err != nil {
		return false, fmt.Errorf("(repo) failed to check if artist is liked by user: %w", err)
	}

	return isLiked, nil
}

func (p *PostgreSQL) UpdateMonthListensPerUser(ctx context.Context) error {
	query := fmt.Sprintf(
		`WITH new_users_by_artist AS (
			SELECT a.id as id, COUNT(DISTINCT tu.user_id) as new_users
			FROM (
				SELECT track_id, user_id
				FROM %s
				WHERE commited_at BETWEEN (current_timestamp - interval '30 day') AND current_timestamp
				GROUP BY track_id, user_id
			) as tu 
				INNER JOIN %s AS at ON tu.track_id = at.track_id
				RIGHT JOIN %s AS a ON at.artist_id = a.id
			GROUP BY a.id, tu.user_id
		)
		UPDATE %s AS a
		SET listens = nla.new_users
		FROM new_users_by_artist AS nla
		WHERE a.id = nla.id;`,
		p.tables.Listens(), p.tables.ArtistsTracks(), p.tables.Artists(), p.tables.Artists())

	if _, err := p.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("(repo) failed to update artist listens per user: %w", err)
	}

	return nil
}
