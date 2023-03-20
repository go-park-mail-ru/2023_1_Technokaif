package track_repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	db "github.com/go-park-mail-ru/2023_1_Technokaif/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// trackPostgres implements TrackRepository
type trackPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewTrackPostgres(db *sqlx.DB, l logger.Logger) track.TrackRepository {
	return &trackPostgres{db: db, logger: l}
}

func (tp *trackPostgres) Insert(track models.Track) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (name, album_id, cover_src, record_src) 
		VALUES ($1, $2, $3, $4);`,
		db.PostgresTables.Tracks)

	if _, err := tp.db.Exec(query, track.Name, track.AlbumID,
		track.CoverSrc, track.RecordSrc); err != nil {

		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (tp *trackPostgres) GetByID(trackID uint32) (models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, cover_src, record_src 
		FROM %s 
		WHERE id = $1;`,
		db.PostgresTables.Tracks)

	var track models.Track
	if err := tp.db.Get(&track, query, trackID); err != nil {
		return models.Track{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return track, nil
}

func (tp *trackPostgres) Update(track models.Track) error {
	query := fmt.Sprintf(
		`UPDATE %s 
		SET name = $1, album_id = $2, cover_src = $3, record_src = $4 
		WHERE id = $5;`,
		db.PostgresTables.Tracks)

	if _, err := tp.db.Exec(query, track.Name, track.AlbumID,
		track.CoverSrc, track.RecordSrc, track.ID); err != nil {

		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (tp *trackPostgres) Delete(trackID uint32) error {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1;`,
		db.PostgresTables.Tracks)

	if _, err := tp.db.Exec(query, trackID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (tp *trackPostgres) GetFeed() ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, cover_src, record_src 
		FROM %s 
		LIMIT 100;`,
		db.PostgresTables.Tracks)

	var tracks []models.Track
	if err := tp.db.Select(&tracks, query); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (tp *trackPostgres) GetByAlbum(albumID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, cover_src, record_src 
		FROM %s
		WHERE album_id = $1;`,
		db.PostgresTables.Tracks)

	var tracks []models.Track
	if err := tp.db.Select(&tracks, query, albumID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (tp *trackPostgres) GetByArtist(artistID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT t.id, t.name, t.album_id, t.cover_src, t.record_src 
		FROM %s t 
			INNER JOIN %s at ON t.id = at.track_id 
		WHERE at.artist_id = $1;`,
		db.PostgresTables.Tracks, db.PostgresTables.ArtistsTracks)

	var tracks []models.Track
	if err := tp.db.Select(&tracks, query, artistID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}

func (tp *trackPostgres) GetLikedByUser(userID uint32) ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT t.id, name, t.album_id, t.cover_src, t.record_src 
		FROM %s t 
			INNER JOIN %s ut ON t.id = ut.track_id 
		WHERE ut.user_id = $1;`,
		db.PostgresTables.Tracks, db.PostgresTables.LikedTracks)

	var tracks []models.Track
	if err := tp.db.Select(&tracks, query, userID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return tracks, nil
}
