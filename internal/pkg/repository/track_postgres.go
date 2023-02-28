package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// TrackPostgres implements Track
type TrackPostgres struct {
	db *sql.DB
}

func NewTrackPostgres(db *sql.DB) *TrackPostgres {
	return &TrackPostgres{db: db}
}

type artistsTracks struct {
	TrackID    int
	TrackName  string
	ArtistID   int
	ArtistName string
}

func (tp *TrackPostgres) GetFeed() ([]models.TrackFeed, error) {
	query := fmt.Sprintf("SELECT t.id, t.name, a.id, a.name "+
		"FROM %s t INNER JOIN %s at ON t.id = at.track_id "+
		"INNER JOIN %s a ON at.artist_id = a.id;",
		TracksTable, ArtistsTracksTable, ArtistsTable)

	rows, err := tp.db.Query(query)
	if err != nil {
		return nil, err
	}

	var m = make(map[int]models.TrackFeed)
	for rows.Next() {
		var at artistsTracks
		if err = rows.Scan(&at.TrackID, &at.TrackName, &at.ArtistID, &at.ArtistName); err != nil {
			return nil, err
		}

		if tf, ok := m[at.TrackID]; !ok {
			m[at.TrackID] = models.TrackFeed{ID: at.TrackID, Name: at.TrackName,
				Artists: []models.ArtistFeed{{ID: at.ArtistID, Name: at.ArtistName}}}
		} else {
			tf.Artists = append(tf.Artists,
				models.ArtistFeed{ID: at.ArtistID, Name: at.ArtistName})
			m[at.TrackID] = tf
		}
	}

	var tracks []models.TrackFeed
	for _, v := range m {
		tracks = append(tracks, v)
	}

	return tracks, nil
}
