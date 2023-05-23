package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate easyjson -no_std_marshalers user_delivery_models.go

// UploadAvatar
const MaxAvatarMemory = 5 << 20
const avatarFormKey = "avatar"

const (
	userNotFound = "no such user"

	userGetServerError          = "can't get user"
	userUpdateInfoServerError   = "can't update user info"
	userAvatarUploadServerError = "can't upload avatar"

	userAvatarUploadInvalidData     = "invalid avatar data"
	userAvatarUploadInvalidDataType = "invalid avatar data type"

	userUpdatedInfoSuccessfully    = "ok"
	userAvatarUploadedSuccessfully = "ok"
)

//easyjson:json
type userUploadAvatarResponse struct {
	Status string `json:"status"`
}

// Update Info
//
//easyjson:json
type userInfoInput struct {
	Email     string      `json:"email" valid:"required,email,maxstringlength(255)"`
	FirstName string      `json:"firstName" valid:"required,runelength(2|20)"`
	LastName  string      `json:"lastName" valid:"required,runelength(2|20)"`
	Sex       models.Sex  `json:"sex" valid:"required,in(F|M|O)"`
	BirthDate models.Date `json:"birthDate" valid:"required,born"`
}

func (ui *userInfoInput) validateAndEscape() error {
	ui.escapeHtml()

	_, err := valid.ValidateStruct(*ui)

	return err
}

func (ui *userInfoInput) escapeHtml() {
	ui.Email = html.EscapeString(ui.Email)
	ui.FirstName = html.EscapeString(ui.FirstName)
	ui.LastName = html.EscapeString(ui.LastName)
}

func (ui *userInfoInput) ToUser(user *models.User) *models.User {
	return &models.User{
		ID:        user.ID,
		Email:     ui.Email,
		FirstName: ui.FirstName,
		LastName:  ui.LastName,
		Sex:       ui.Sex,
		BirthDate: ui.BirthDate,
		AvatarSrc: user.AvatarSrc,
	}
}

//easyjson:json
type userChangeInfoResponse struct {
	Status string `json:"status"`
}
