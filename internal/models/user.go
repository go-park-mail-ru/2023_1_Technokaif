package models

import (
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
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
	ID        uint32 `json:"-" 		   valid:"-"`
	Version   uint32 `json:"-" 		   valid:"-"`
	Username  string `json:"username"  valid:"required,runelength(4|20)"`
	Email     string `json:"email"     valid:"required,email"`
	Password  string `json:"password"  valid:"required,runelength(8|30),passwordcheck"`
	Salt      string `json:"-" 		   valid:"-"`
	FirstName string `json:"firstName" valid:"required,runelength(2|20)"`
	LastName  string `json:"lastName"  valid:"required,runelength(2|20)"`
	Sex       Sex    `json:"sex"       valid:"required,in(F|M|O)"`
	BirhDate  Date   `json:"birthDate" valid:"required,born"`
	// AvatarSrc string `json:"-" valid:"-"`
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	var err error
	d.Time, err = time.Parse("2006-01-02", s) // RFC 3339

	return err
}

func (u *User) DeliveryValidate() error {
	_, err := valid.ValidateStruct(u)
	return err
}
