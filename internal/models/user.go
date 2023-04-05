package models

import (
	"strings"
	"time"
	"html"

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
	ID        uint32 `json:"-" valid:"-" db:"id"`
	Version   uint32 `json:"-" valid:"-" db:"version"`
	Username  string `json:"username" valid:"required,runelength(4|20)" db:"username"`
	Email     string `json:"email" valid:"required,email,maxstringlength(255)" db:"email"`
	Password  string `json:"password" valid:"required,runelength(8|30),passwordcheck" db:"password_hash"`
	Salt      string `json:"-" valid:"-" db:"salt"`
	FirstName string `json:"firstName" valid:"required,runelength(2|20)" db:"first_name"`
	LastName  string `json:"lastName" valid:"required,runelength(2|20)" db:"last_name"`
	Sex       Sex    `json:"sex" valid:"required,in(F|M|O)" db:"sex"`
	BirthDate Date   `json:"birthDate" valid:"required,born" db:"birth_date"`
	AvatarSrc string `json:"avatarSrc" valid:"-" db:"avatar_src"`
}

type UserTransfer struct {
	ID        uint32 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Sex       Sex    `json:"sex"`
	BirhDate  Date   `json:"birthDate"`
	AvatarSrc string `json:"avatarSrc,omitempty"`
}

type ContextKeyUserType struct{}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	var err error
	d.Time, err = time.Parse("2006-01-02", s) // RFC 3339

	return err
}

func (u *User) DeliveryValidate() error {
	u.Username = html.EscapeString(u.Username)
	u.Email = html.EscapeString(u.Email)
	u.Password = html.EscapeString(u.Password)
	u.FirstName = html.EscapeString(u.FirstName)
	u.LastName = html.EscapeString(u.LastName)

	_, err := valid.ValidateStruct(u)
	return err
}

func UserTransferFromUser(user User) UserTransfer {
	return UserTransfer{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Sex:       user.Sex,
		BirhDate:  user.BirthDate,
		AvatarSrc: user.AvatarSrc,
	}
}
