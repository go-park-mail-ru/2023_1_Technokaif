package models

type Album struct {
	ID          uint32  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	CoverSrc    string  `db:"cover_src"`
}
