package common

import (
	"time"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	commonProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/common/proto/generated"
)

func UserToProto(user *models.User) *commonProto.UserResponse {
	return &commonProto.UserResponse{
		Id:           user.ID,
		Version:      user.Version,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Sex:          string(user.Sex),
		AvatarSrc:    user.AvatarSrc,
		BirthDate:    user.BirthDate.Format("2006-01-02"),
	}
}

func ProtoToUser(userProto *commonProto.UserResponse) (*models.User, error) {
	time, err := time.Parse("2006-01-02", userProto.BirthDate)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        userProto.Id,
		Version:   userProto.Version,
		Username:  userProto.Username,
		Email:     userProto.Email,
		Password:  userProto.PasswordHash,
		FirstName: userProto.FirstName,
		LastName:  userProto.LastName,
		Sex:       models.Sex(userProto.Sex),
		AvatarSrc: userProto.AvatarSrc,
		BirthDate: models.Date{Time: time},
	}

	return user, nil
}
