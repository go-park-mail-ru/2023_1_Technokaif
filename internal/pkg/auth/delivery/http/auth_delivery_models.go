package http

import (
	"fmt"
	"html"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate easyjson -no_std_marshalers auth_delivery_models.go

// Response messages
const (
	userNotFound = "no such user"

	userAlreadyExists = "user already exists"
	passwordMismatch  = "incorrect password"
	invalidToken      = "invalidToken"
	userForbidden     = "forbidden"

	userSignUpServerError    = "can't sign up user"
	userLoginServerError     = "can't login user"
	userLogoutServerError    = "can't log out user"
	userGetServerError       = "can't get user"
	userChangePasswordError  = "can't change password"
	tokenGenerateServerError = "can't generate new token"

	userLogedOutSuccessfully        = "ok"
	userChangedPasswordSuccessfully = "ok"
)

// Signup
//
//easyjson:json
type signUpInput struct {
	Username  string      `json:"username" valid:"required,runelength(4|20)"`
	Email     string      `json:"email" valid:"required,email,maxstringlength(255)"`
	Password  string      `json:"password" valid:"required,runelength(8|30),passwordcheck"`
	FirstName string      `json:"firstName" valid:"required,runelength(2|20)"`
	LastName  string      `json:"lastName" valid:"required,runelength(2|20)"`
	BirthDate models.Date `json:"birthDate" valid:"required,born"`
	AvatarSrc string      `json:"avatarSrc" valid:"-"`
}

func (sui *signUpInput) ToUser() models.User {
	return models.User{
		Username:  sui.Username,
		Email:     sui.Email,
		Password:  sui.Password,
		FirstName: sui.FirstName,
		LastName:  sui.LastName,
		BirthDate: sui.BirthDate,
		AvatarSrc: sui.AvatarSrc,
	}
}

//easyjson:json
type signUpResponse struct {
	ID uint32 `json:"id"`
}

// Login
//
//easyjson:json
type loginInput struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (li *loginInput) validateAndEscape() error {
	li.escapeHtml()

	_, err := valid.ValidateStruct(*li)

	return err
}

func (li *loginInput) escapeHtml() {
	li.Username = html.EscapeString(li.Username)
	li.Password = html.EscapeString(li.Password)
}

//easyjson:json
type loginResponse struct {
	UserID uint32 `json:"id"`
}

// Logout
//
//easyjson:json
type logoutResponse struct {
	Status string `json:"status"`
}

// ChangePassword
//
//easyjson:json
type changePassInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword" valid:"required,runelength(8|30),passwordcheck"`
}

func (c *changePassInput) validate() error {
	c.NewPassword = html.EscapeString(c.NewPassword)

	_, err := valid.ValidateStruct(*c)

	return err
}

//easyjson:json
type changePassResponse struct {
	Status string `json:"status"`
}

//easyjson:json
type isAuthenticatedResponse struct {
	Authenticated bool `json:"auth"`
}

type signUpValidateErrors struct {
	Err    error
	Fields []string
}

func (e *signUpValidateErrors) HttpErrorResponce() string {
	if len(e.Fields) == 1 {
		return fmt.Sprintf("incorrect field: %s", e.Fields[0])
	}
	return fmt.Sprintf("incorrect fields: %s", strings.Join(e.Fields, ", "))
}

func (e *signUpValidateErrors) Error() string {
	return e.Err.Error()
}

func signUpInputAuthDeliveryValidate(user *signUpInput) error {
	user.Username = html.EscapeString(user.Username)
	user.Email = html.EscapeString(user.Email)
	user.Password = html.EscapeString(user.Password)
	user.FirstName = html.EscapeString(user.FirstName)
	user.LastName = html.EscapeString(user.LastName)

	_, err := valid.ValidateStruct(*user)
	if errsValidate, ok := err.(valid.Errors); ok {
		var fields []string
		for _, err := range errsValidate {
			if errValidate, ok := err.(valid.Error); ok {
				fields = append(fields, errValidate.Name)
			}
		}

		if len(fields) > 0 {
			return &signUpValidateErrors{
				Fields: fields,
				Err:    err,
			}
		}
	}

	return err
}
