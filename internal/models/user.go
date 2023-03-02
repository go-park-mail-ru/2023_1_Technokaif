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
	ID        	uint   `json:"-" 		 valid:"-"`
	Version 	uint   `json:"-" 		 valid:"-"`
	Username  	string `json:"username"  valid:"required,runelength(4|20)"`
	Email     	string `json:"email"     valid:"required,email"`
	Password  	string `json:"password"  valid:"required,runelength(8|30),passwordcheck"`
	FirstName 	string `json:"firstName" valid:"required,runelength(2|20)"`
	LastName  	string `json:"lastName"  valid:"required,runelength(2|20)"`
	Sex       	Sex    `json:"sex"       valid:"required,in(F|M|O)"`
	BirhDate  	Date   `json:"birthDate" valid:"required,born"`
	// AvatarSrc string `json:"-" valid:"-"`
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s) // RFC 3339
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}
