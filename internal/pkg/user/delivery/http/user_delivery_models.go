package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// UploadAvatar
const MaxAvatarMemory = 1 << 20
const avatarForm = "avatar"

type userUploadAvatarResponse struct {
	Status string `json:"status"`
}

// Update Info
type userInfoInput struct {
	Email     string      `json:"email" valid:"required,email,maxstringlength(255)"`
	FirstName string      `json:"firstName" valid:"required,runelength(2|20)"`
	LastName  string      `json:"lastName" valid:"required,runelength(2|20)"`
	Sex       models.Sex  `json:"sex" valid:"required,in(F|M|O)"`
	BirthDate models.Date `json:"birthDate" valid:"required,born"`
}

func (ui *userInfoInput) validate() error {
	ui.Email = html.EscapeString(ui.Email)
	ui.FirstName = html.EscapeString(ui.FirstName)
	ui.LastName = html.EscapeString(ui.LastName)

	_, err := valid.ValidateStruct(*ui)

	return err
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

type userChangeInfoResponse struct {
	Status string `json:"status"`
}
