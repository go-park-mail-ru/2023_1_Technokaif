package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type AlbumPostgres struct {
	db *sql.DB
}

func NewAlbumPostgres(db *sql.DB) *AlbumPostgres {
	return &AlbumPostgres{db: db}
}

func (tp *AlbumPostgres) GetFeed() ([]models.AlbumFeed, error) {
	query := fmt.Sprintf("SELECT al.name, a.name "+
		"FROM %s al INNER JOIN %s a ON al.artist_id = a.id;",
		AlbumsTable, ArtistsTable)

	rows, err := tp.db.Query(query)
	if err != nil {
		return nil, err
	}

	var albums []models.AlbumFeed
	for rows.Next() {
		var album models.AlbumFeed
		if err = rows.Scan(&album.Name, &album.ArtistName); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}
