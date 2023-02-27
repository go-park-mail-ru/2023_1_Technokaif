package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type TrackPostgres struct {
	db *sql.DB
}

func NewTrackPostgres(db *sql.DB) *TrackPostgres {
	return &TrackPostgres{db: db}
}

func (tp *TrackPostgres) GetFeed() ([]models.TrackFeed, error) {
	query := fmt.Sprintf("SELECT t.name, a.name "+
		"FROM %s t INNER JOIN %s a ON t.artist_id = a.id",
		TracksTable, ArtistsTable)

	rows, err := tp.db.Query(query)
	if err != nil {
		fmt.Println("Err query")
		return nil, err
	}

	var tracks []models.TrackFeed
	for rows.Next() {
		var track models.TrackFeed
		if err = rows.Scan(&track.Name, &track.ArtistName); err != nil {
			fmt.Println("Err scan")
			return nil, err
		}
		tracks = append(tracks, track)
	}

	fmt.Println("Success")

	return tracks, nil
}
