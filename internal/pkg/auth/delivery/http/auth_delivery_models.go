package http

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
)

// Signup
type signUpResponse struct {
	ID uint32 `json:"id"`
}

// Login
type loginInput struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (li *loginInput) validate() error {
	sanitizer := bluemonday.StrictPolicy()
	li.Username = sanitizer.Sanitize(li.Username)
	li.Password = sanitizer.Sanitize(li.Password)
	
	_, err := valid.ValidateStruct(li)

	return err
}

type loginResponse struct {
	JWT string `json:"jwt"`
}

// Logout
type logoutResponse struct {
	Status string `json:"status"`
}
