package models

import (
	"fmt"
	"strings"
	"time"
)

type Sex string

const (
	Male   Sex = "M"
	Female Sex = "F"
	Other  Sex = "O"
)


// User implements information about app's users
type Date struct {
	time.Time
}
type User struct {
	ID        uint   `json:"-" valid:"-"`
	Username  string `json:"username"  valid:"required,runelength(4|20)"`
	Email     string `json:"email"     valid:"required,email"`
	Password  string `json:"password"  valid:"required,runelength(8|20),passwordcheck"`
	FirstName string `json:"firstName" valid:"required,runelength(2|20)"`
	LastName  string `json:"lastName"  valid:"required,runelength(2|20)"`
	Sex       Sex    `json:"sex"       valid:"required,in(F|M|O)"`
	BirhDate  Date	 `json:"birthDate" valid:"required,born"`
}

func (d *Date) UnmarshalJSON(b []byte) error {
    s := strings.Trim(string(b), "\"")
    t, err := time.Parse("2006-01-02", s)  // RFC 3339
    if err != nil {
		fmt.Println("UnmarshalJSON")
        return err
    }

    d.Time = t
    return nil
}
