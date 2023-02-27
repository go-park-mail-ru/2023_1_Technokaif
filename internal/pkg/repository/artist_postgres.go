package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type ArtistPostgres struct {
	db *sql.DB
}

func NewArtistPostgres(db *sql.DB) *ArtistPostgres {
	return &ArtistPostgres{db: db}
}

func (tp *ArtistPostgres) GetFeed() ([]models.ArtistFeed, error) {
	query := fmt.Sprintf("SELECT name FROM %s;", ArtistsTable)

	rows, err := tp.db.Query(query)
	if err != nil {
		return nil, err
	}

	var artists []models.ArtistFeed
	for rows.Next() {
		var artist models.ArtistFeed
		if err = rows.Scan(&artist.Name); err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}

	return artists, nil
}
