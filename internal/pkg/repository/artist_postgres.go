package repository

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// ArtistPostgres implements Artist
type ArtistPostgres struct {
	db     *sql.DB
	logger logger.Logger
}

func NewArtistPostgres(db *sql.DB, l logger.Logger) *ArtistPostgres {
	return &ArtistPostgres{db: db, logger: l}
}

func (tp *ArtistPostgres) GetFeed() ([]models.ArtistFeed, error) {
	query := fmt.Sprintf("SELECT id, name, avatar_src FROM %s;", artistsTable)

	rows, err := tp.db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "(Repo) failed to make query")
	}
	defer rows.Close()

	var artists []models.ArtistFeed
	for rows.Next() {
		var artist models.ArtistFeed
		if err = rows.Scan(&artist.ID, &artist.Name, &artist.AvatarSrc); err != nil {
			return nil, errors.Wrap(err, "(Repo) failed to scan from query")
		}
		artists = append(artists, artist)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "(Repo) failed to scan from query")
	}

	return artists, nil
}
