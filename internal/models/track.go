package models

type Track struct {
	ID            uint32  `db:"id"`
	Name          string  `db:"name"`
	AlbumID       *uint32 `db:"album_id"`
	AlbumPosition *uint32 `db:"album_position"`
	CoverSrc      string  `db:"cover_src"`
	RecordSrc     string  `db:"record_src"`
	Duration      uint32  `db:"duration"`
	Listens       uint32  `db:"listens"`
}
