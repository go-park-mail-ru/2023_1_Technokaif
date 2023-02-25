package repository

import (
	"database/sql"
)

type AlbumPostgres struct {
	db *sql.DB
}

func NewAlbumPostgres(db *sql.DB) *AlbumPostgres {
	return &AlbumPostgres{db: db}
}