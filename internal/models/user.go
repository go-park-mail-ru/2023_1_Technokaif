package models

import (
	"strings"
	"time"
)

//go:generate easyjson -no_std_marshalers user.go

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
	BirthDate Date   `db:"birth_date"`
	AvatarSrc string `db:"avatar_src"`
}

//easyjson:json
type UserTransfer struct {
	ID        uint32 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate Date   `json:"birthDate,omitempty"`
	AvatarSrc string `json:"avatarSrc,omitempty"`
}

//easyjson:json
type UserTransfers []UserTransfer

func UserTransferFromEntry(user User) UserTransfer {
	return UserTransfer{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		AvatarSrc: user.AvatarSrc,
	}
}

func UserTransferFromList(users []User) []UserTransfer {
	userTransfers := make([]UserTransfer, 0, len(users))
	for _, u := range users {
		userTransfers = append(userTransfers, UserTransferFromEntry(u))
	}

	return userTransfers
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	var err error
	d.Time, err = time.Parse("2006-01-02", s) // RFC 3339

	return err
}
