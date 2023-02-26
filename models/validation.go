package models

import (
	"time"
	"unicode"

	valid "github.com/asaskevich/govalidator"
)

func init() { // For validation
	valid.CustomTypeTagMap.Set("born", func(date interface{}, context interface{}) bool {
		if _, ok := context.(User); !ok {
			return false
		}

		d, ok := date.(Date)
		if !ok {
			return false
		}

		return d.Time.Before(time.Now())
	})

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
}

func (u *User) DeliveryValidate() bool {
	isValid, err := valid.ValidateStruct(u)
	if err != nil || !isValid {
		return false
	}

	return true
}
