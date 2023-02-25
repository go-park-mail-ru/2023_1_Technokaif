package models

import (
	"unicode"

	valid "github.com/asaskevich/govalidator"
)

type Sex string

const (
	Male   Sex = "M"
	Female Sex = "F"
	Other  Sex = "O"
)

// User implements information about app's users
type User struct {
	ID        uint   `json:"-" valid:"-"`
	Username  string `json:"username"  valid:"runelength(4|20)"`
	Email     string `json:"email"     valid:"email"`
	Password  string `json:"password"  valid:"runelength(8|20),passwordcheck"`
	FirstName string `json:"firstName" valid:"runelength(2|20)"`
	LastName  string `json:"lastName"  valid:"runelength(2|20)"`
	Sex       Sex    `json:"sex"       valid:"in(F|M|O)"`
}

func (u *User) DeliveryValidate() bool {
	valid.TagMap["passwordcheck"] = valid.Validator(func(password string) bool {
		hasLowLetters := false
		hasUpperLetters := false
		hasDigits := false

		for _, c := range password {
			switch {
			case unicode.IsNumber(c):
				hasDigits = true
			case unicode.IsUpper(c):
				hasUpperLetters = true
			case unicode.IsLower(c):
				hasLowLetters = true
			}
		}

		return hasLowLetters && hasUpperLetters && hasDigits
	})

	valid.SetFieldsRequiredByDefault(true)
	isValid, err := valid.ValidateStruct(u)
	if err != nil || !isValid {
		return false
	}

	return true
}
