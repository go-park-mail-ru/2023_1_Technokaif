package models

type Artist struct {
	ID        uint32  `db:"id"`
	UserID    *uint32 `db:"user_id"`
	Name      string  `db:"name"`
	AvatarSrc string  `db:"avatar_src"`
}
