package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

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
type signUpResponse struct {
	ID uint32 `json:"id"`
}

// Login
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

type loginResponse struct {
	UserID uint32 `json:"id"`
}

// Logout
type logoutResponse struct {
	Status string `json:"status"`
}

// ChangePassword
type changePassInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword" valid:"required,runelength(8|30),passwordcheck"`
}

func (c *changePassInput) validate() error {
	c.NewPassword = html.EscapeString(c.NewPassword)

	_, err := valid.ValidateStruct(*c)

	return err
}

type changePassResponse struct {
	Status string `json:"status"`
}

type isAuthenticatedResponse struct {
	Authenticated bool `json:"auth"`
}

func userAuthDeliveryValidate(user *models.User) error {
	user.Username = html.EscapeString(user.Username)
	user.Email = html.EscapeString(user.Email)
	user.Password = html.EscapeString(user.Password)
	user.FirstName = html.EscapeString(user.FirstName)
	user.LastName = html.EscapeString(user.LastName)

	_, err := valid.ValidateStruct(*user)
	return err
}
