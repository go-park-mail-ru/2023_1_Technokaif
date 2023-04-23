package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// PostgreSQL implements album.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables playlist.Tables
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, t playlist.Tables, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
		logger: l,
	}
}

func checkTransaction(tx *sql.Tx, repoError *error) {
	if *repoError != nil {
		if err := tx.Rollback(); err != nil {
			*repoError = fmt.Errorf("(repo) failed to Rollback: %w: %w", err, *repoError)
		}

	} else {
		if err := tx.Commit(); err != nil {
			*repoError = fmt.Errorf("(repo) failed to Commit: %w: %w", err, *repoError)
		}
	}
}

const errorAlreadyExists = "unique_violation"

func (p *PostgreSQL) Insert(playlist models.Playlist, usersID []uint32) (_ uint32, repoErr error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer checkTransaction(tx, &repoErr)

	insertAlbumQuery := fmt.Sprintf(
		`INSERT INTO %s (name, description, cover_src)
		VALUES ($1, $2, $3) RETURNING id;`,
		p.tables.Playlists())

	var playlistID uint32
	row := tx.QueryRow(insertAlbumQuery, playlist.Name, playlist.Description, playlist.CoverSrc)
	if err := row.Scan(&playlistID); err != nil {
		return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	insertPlaylistUsersQuery := fmt.Sprintf(
		`INSERT INTO %s (user_id, playlist_id)
		VALUES ($1, $2);`,
		p.tables.UsersPlaylists())

	for _, userID := range usersID {
		if _, err := tx.Exec(insertPlaylistUsersQuery, userID, playlistID); err != nil {
			return 0, fmt.Errorf("(repo) failed to exec query: %w", err)
		}
	}

	return playlistID, nil
}

func (p *PostgreSQL) GetByID(playlistID uint32) (*models.Playlist, error) {
	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src 
		FROM %s 
		WHERE id = $1;`,
		p.tables.Playlists())

	var playlist models.Playlist
	if err := p.db.Get(&playlist, query, playlistID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchPlaylistError{PlaylistID: playlistID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &playlist, nil
}

func (p *PostgreSQL) UpdateWithMembers(pl models.Playlist, usersID []uint32) (repoErr error) {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer checkTransaction(tx, &repoErr)

	updatePlaylistQuery := fmt.Sprintf(
		`UPDATE %s
		SET name = $2,
			description = $3,
			cover_src = $4
		WHERE id = $1;`,
		p.tables.Playlists())

	if _, err := p.db.Exec(updatePlaylistQuery, pl.ID, pl.Name, pl.Description, pl.CoverSrc); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	insertPlaylistUsersQuery := fmt.Sprintf(
		`INSERT INTO %s (user_id, playlist_id)
		VALUES ($1, $2);`,
		p.tables.UsersPlaylists())

	for _, userID := range usersID {
		if _, err := tx.Exec(insertPlaylistUsersQuery, userID, pl.ID); err != nil {
			return fmt.Errorf("(repo) failed to exec query: %w", err)
		}
	}

	return nil
}

func (p *PostgreSQL) Update(pl models.Playlist) error {
	updatePlaylistQuery := fmt.Sprintf(
		`UPDATE %s
		SET name = $2,
			description = $3,
			cover_src = $4
		WHERE id = $1;`,
		p.tables.Playlists())

	if _, err := p.db.Exec(updatePlaylistQuery, pl.ID, pl.Name, pl.Description, pl.CoverSrc); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) DeleteByID(playlistID uint32) error {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE id = $1;`,
		p.tables.Playlists())

	resExec, err := p.db.Exec(query, playlistID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return fmt.Errorf("(repo) failed to check RowsAffected: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("(repo): %w", &models.NoSuchPlaylistError{PlaylistID: playlistID})
	}

	return nil
}

func (p *PostgreSQL) AddTrack(trackID, playlistID uint32) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (track_id, playlist_id)
		VALUES ($1, $2);`,
		p.tables.PlaylistsTracks())

	if _, err := p.db.Exec(query, trackID, playlistID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("(repo) %w: %w", &models.NoSuchPlaylistError{PlaylistID: playlistID}, err)
		}

		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == errorAlreadyExists {
				return fmt.Errorf("(repo) entry already exists: %w", pqerr)
			}
		}

		return fmt.Errorf("(repo) failed to insert: %w", err)
	}

	return nil
}

func (p *PostgreSQL) DeleteTrack(trackID, playlistID uint32) error {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE track_id = $1 AND playlist_id = $2;`,
		p.tables.PlaylistsTracks())

	resExec, err := p.db.Exec(query, trackID, playlistID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	deleted, err := resExec.RowsAffected()
	if err != nil {
		return fmt.Errorf("(repo) failed to check RowsAffected: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("(repo) no such track or playlist")
	}

	return nil
}

func (p *PostgreSQL) GetFeed() ([]models.Playlist, error) {
	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src  
		FROM %s 
		LIMIT 100;`,
		p.tables.Playlists())

	var playlists []models.Playlist
	if err := p.db.Select(&playlists, query); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return playlists, nil
}

func (p *PostgreSQL) GetByUser(userID uint32) ([]models.Playlist, error) {
	query := fmt.Sprintf(
		`SELECT p.id, p.name, p.description, p.cover_src 
		FROM %s p
			INNER JOIN %s up ON p.id = up.playlist_id
		WHERE up.user_id = $1;`,
		p.tables.Playlists(), p.tables.UsersPlaylists())

	var playlists []models.Playlist
	if err := p.db.Select(&playlists, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return playlists, nil
}

func (p *PostgreSQL) GetLikedByUser(userID uint32) ([]models.Playlist, error) {
	query := fmt.Sprintf(
		`SELECT p.id, p.name, p.description, p.cover_src
		FROM %s p 
			INNER JOIN %s ua ON p.id = up.playlist_id 
		WHERE up.user_id = $1;`,
		p.tables.Playlists(), p.tables.LikedPlaylists())

	var playlists []models.Playlist
	if err := p.db.Select(&playlists, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return playlists, nil
}

func (p *PostgreSQL) InsertLike(playlistID, userID uint32) (bool, error) {
	insertLikeQuery := fmt.Sprintf(
		`INSERT INTO %s (playlist_id, user_id) 
		VALUES ($1, $2);`,
		p.tables.LikedPlaylists())

	if _, err := p.db.Exec(insertLikeQuery, playlistID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("(repo) %w: %w", &models.NoSuchPlaylistError{PlaylistID: playlistID}, err)
		}

		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == errorAlreadyExists {
				return false, nil
			}
		}

		return false, fmt.Errorf("(repo) failed to insert: %w", err)
	}

	return true, nil
}

func (p *PostgreSQL) DeleteLike(playlistID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`DELETE
		FROM %s
		WHERE playlist_id = $1 AND user_id = $2;`,
		p.tables.LikedPlaylists())

	resExec, err := p.db.Exec(query, playlistID, userID)
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

func (p *PostgreSQL) IsLiked(playlistID, userID uint32) (bool, error) {
	query := fmt.Sprintf(
		`SELECT CASE WHEN 
			EXISTS(SELECT *
				FROM %s
				WHERE playlist_id = $1 AND user_id = $2
			) THEN TRUE ELSE FALSE END;`,
		p.tables.LikedPlaylists())

	var isLiked bool
	err := p.db.Get(&isLiked, query, playlistID, userID)
	if err != nil {
		return false, fmt.Errorf("(repo) failed to check if playlist is liked by user: %w", err)
	}

	return isLiked, nil
}
