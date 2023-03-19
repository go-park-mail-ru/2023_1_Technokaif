package track_repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/go-park-mail-ru/2023_1_Technokaif/init/db"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// TrackPostgres implements Track
type TrackPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewTrackPostgres(db *sqlx.DB, l logger.Logger) *TrackPostgres {
	return &TrackPostgres{db: db, logger: l}
}

func (tp *TrackPostgres) GetByID(trackID uint32) (models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, cover_src, record_src 
		FROM %s 
		WHERE id = $1;`,
		db.PostgresTables.Tracks)

	var track models.Track
	if err := tp.db.Get(&track, query, trackID); err != nil {
		return models.Track{}, fmt.Errorf("failed to exec query: %w", err)
	}

	return track, nil
}

func (tp *TrackPostgres) GetFeed() ([]models.Track, error) {
	query := fmt.Sprintf(
		`SELECT id, name, album_id, cover_src, record_src 
		FROM %s 
		ORDER BY id;`,
		db.PostgresTables.Tracks)

	var tracks []models.Track
	if err := tp.db.Select(&tracks, query); err != nil {
		return nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return tracks, nil
}

func (tp *TrackPostgres) GetByArtist(artistID uint32) ([]models.Track, error) {

}

func (tp *TrackPostgres) GetByAlbum(albumID uint32) ([]models.Track, error) {

}

func (tp *TrackPostgres) GetByUser(userID uint32) ([]models.Track, error) {

}
