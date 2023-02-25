package repository

import (
	"database/sql"
)

type TrackPostgres struct {
	db *sql.DB
}

func NewTrackPostgres(db *sql.DB) *TrackPostgres {
	return &TrackPostgres{db: db}
}