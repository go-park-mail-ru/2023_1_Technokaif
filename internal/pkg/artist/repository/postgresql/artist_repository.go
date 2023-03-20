package artist_repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	db "github.com/go-park-mail-ru/2023_1_Technokaif/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// artistPostgres implements ArtistRepository
type artistPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewArtistPostgres(db *sqlx.DB, l logger.Logger) artist.ArtistRepository {
	return &artistPostgres{db: db, logger: l}
}

func (ap *artistPostgres) Insert(artist models.Artist) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (name, avatar_src) 
		VALUES ($1, $2);`,
		db.PostgresTables.Artists)

	if _, err := ap.db.Exec(query, artist.Name, artist.AvatarSrc); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (ap *artistPostgres) GetByID(artistID uint32) (models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT id, name, avatar_src 
		FROM %s 
		WHERE id = $1;`,
		db.PostgresTables.Artists)

	var artist models.Artist
	if err := ap.db.Get(&artist, query, artistID); err != nil {
		return models.Artist{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artist, nil
}

func (ap *artistPostgres) Update(artist models.Artist) error {
	query := fmt.Sprintf(
		`UPDATE %s 
		SET name = $1, avatar_src = $2 
		WHERE id = $3;`,
		db.PostgresTables.Artists)

	if _, err := ap.db.Exec(query, artist.Name, artist.AvatarSrc, artist.ID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (ap *artistPostgres) DeleteByID(artistID uint32) error {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1;`,
		db.PostgresTables.Artists)

	if _, err := ap.db.Exec(query, artistID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (ap *artistPostgres) GetFeed() ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT id, name, avatar_src  
		FROM %s 
		LIMIT 100;`,
		db.PostgresTables.Artists)

	var artists []models.Artist
	if err := ap.db.Select(&artists, query); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (ap *artistPostgres) GetByAlbum(albumID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src 
		FROM %s a 
			INNER JOIN %s aa ON a.id = aa.artist_id 
		WHERE aa.album_id = $1;`,
		db.PostgresTables.Artists, db.PostgresTables.ArtistsAlbums)

	var artists []models.Artist
	if err := ap.db.Select(&artists, query, albumID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (ap *artistPostgres) GetByTrack(trackID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src 
		FROM %s a 
			INNER JOIN %s at ON a.id = at.artist_id 
		WHERE at.track_id = $1;`,
		db.PostgresTables.Artists, db.PostgresTables.ArtistsTracks)

	var artists []models.Artist
	if err := ap.db.Select(&artists, query, trackID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}

func (ap *artistPostgres) GetLikedByUser(userID uint32) ([]models.Artist, error) {
	query := fmt.Sprintf(
		`SELECT a.id, a.name, a.avatar_src, 
		FROM %s a 
			INNER JOIN %s ua ON a.id = ua.artist_id 
		WHERE ua.user_id = $1;`,
		db.PostgresTables.Artists, db.PostgresTables.LikedArtists)

	var artists []models.Artist
	if err := ap.db.Select(&artists, query, userID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return artists, nil
}
