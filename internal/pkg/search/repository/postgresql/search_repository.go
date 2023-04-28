package postgresql

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search"
	"github.com/jmoiron/sqlx"
)

// PostgreSQL implements artist.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables search.Tables
}

func NewPostgreSQL(db *sqlx.DB, t search.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

func (p *PostgreSQL) FullTextSearchAlbums(
	ctx context.Context, ftsQuery string, limit uint32) ([]models.Album, error) {

	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src
		FROM %s
		WHERE to_tsvector(lang, name) @@ plainto_tsquery($1)
			OR name LIKE '%%$1%%'
		ORDER BY ts_rank(to_tsvector(lang, name), plainto_tsquery($1)) DESC
		LIMIT $2;`,
		p.tables.Albums(),
	)

	var albums []models.Album
	if err := p.db.SelectContext(ctx, &albums, query, ftsQuery, limit); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return albums, nil
}

func (p *PostgreSQL) FullTextSearchArtists(
	ctx context.Context, ftsQuery string, limit uint32) ([]models.Artist, error) {

	query := fmt.Sprintf(
		`SELECT id, name, avatar_src
		FROM %s
		WHERE to_tsvector(lang, name) @@ plainto_tsquery($1)
			OR name LIKE '%%$1%%'
		ORDER BY ts_rank(to_tsvector(lang, name), plainto_tsquery($1)) DESC
		LIMIT $2;`,
		p.tables.Artists(),
	)

	var artists []models.Artist
	if err := p.db.SelectContext(ctx, &artists, query, ftsQuery, limit); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (p *PostgreSQL) FullTextSearchTracks(
	ctx context.Context, ftsQuery string, limit uint32) ([]models.Track, error) {

	query := fmt.Sprintf(
		`SELECT id, name, cover_src, record_src, duration, listens
		FROM %s
		WHERE to_tsvector(lang, name) @@ plainto_tsquery($1)
			OR name LIKE '%%$1%%'
		ORDER BY ts_rank(to_tsvector(lang, name), plainto_tsquery($1)) DESC
		LIMIT $2;`,
		p.tables.Tracks(),
	)

	var tracks []models.Track
	if err := p.db.SelectContext(ctx, &tracks, query, ftsQuery, limit); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (p *PostgreSQL) FullTextSearchPlaylists(
	ctx context.Context, ftsQuery string, limit uint32) ([]models.Playlist, error) {

	query := fmt.Sprintf(
		`SELECT id, name, description, cover_src
		FROM %s
		WHERE to_tsvector(lang, name) @@ plainto_tsquery($1)
			OR name LIKE '%%$1%%'
		ORDER BY ts_rank(to_tsvector(lang, name), plainto_tsquery($1)) DESC
		LIMIT $2;`,
		p.tables.Playlists(),
	)

	var playlists []models.Playlist
	if err := p.db.SelectContext(ctx, &playlists, query, ftsQuery, limit); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return playlists, nil
}
