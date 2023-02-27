package repository

import (
	"database/sql"
)

type ArtistPostgres struct {
	db *sql.DB
}

func NewArtistPostgres(db *sql.DB) *ArtistPostgres {
	return &ArtistPostgres{db: db}
}
