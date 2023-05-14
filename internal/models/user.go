package models

import (
	"strings"
	"time"
)

type Sex string

const (
	Male   Sex = "M"
	Female Sex = "F"
	Other  Sex = "O"
)

type Date struct {
	time.Time
}

// User implements information about app's users
type User struct {
	ID        uint32 `db:"id"`
	Version   uint32 `db:"version"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	Password  string `db:"password_hash"`
	Salt      string `db:"salt"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Sex       Sex    `db:"sex"`
	BirthDate Date   `db:"birth_date"`
	AvatarSrc string `db:"avatar_src"`
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	var err error
	d.Time, err = time.Parse("2006-01-02", s) // RFC 3339

	return err
}
