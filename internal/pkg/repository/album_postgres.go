package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// AlbumPostgres implements Album
type AlbumPostgres struct {
	db     *sql.DB
	logger logger.Logger
}

func NewAlbumPostgres(db *sql.DB, l logger.Logger) *AlbumPostgres {
	return &AlbumPostgres{db: db, logger: l}
}

type artistsAlbums struct {
	AlbumID     int
	AlbumName   string
	AlbumCover  string
	ArtistID    int
	ArtistName  string
	Description string
}

func (ap *AlbumPostgres) GetFeed() ([]models.AlbumFeed, error) {
	query := fmt.Sprintf(
		"SELECT al.id, al.name, al.cover_src, ar.id, ar.name, al.description "+
			"FROM %s al INNER JOIN %s aa ON al.id = aa.album_id "+
			"INNER JOIN %s ar ON aa.artist_id = ar.id;",
		AlbumsTable, ArtistsAlbumsTable, ArtistsTable)

	rows, err := ap.db.Query(query)
	if err != nil {
		ap.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var m = make(map[int]models.AlbumFeed)
	for rows.Next() {
		var aa artistsAlbums
		if err = rows.Scan(&aa.AlbumID, &aa.AlbumName, &aa.AlbumCover,
			&aa.ArtistID, &aa.ArtistName, &aa.Description); err != nil {
			ap.logger.Error(err.Error())
			return nil, err
		}

		if af, ok := m[aa.AlbumID]; !ok {
			m[aa.AlbumID] = models.AlbumFeed{
				ID:          aa.AlbumID,
				Name:        aa.AlbumName,
				CoverSrc:    aa.AlbumCover,
				Artists:     []models.ArtistFeed{{ID: aa.ArtistID, Name: aa.ArtistName}},
				Description: aa.Description}
		} else {
			af.Artists = append(af.Artists,
				models.ArtistFeed{ID: aa.ArtistID, Name: aa.ArtistName})
			m[aa.AlbumID] = af
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var albums []models.AlbumFeed
	for _, v := range m {
		albums = append(albums, v)
	}

	return albums, nil
}
