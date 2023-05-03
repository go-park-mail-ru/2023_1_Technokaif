package common

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	commonProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/common/proto/generated"
)

func UserToProto(user models.User) *commonProto.UserResponse {
	return &commonProto.UserResponse{
		Id:           user.ID,
		Version:      user.Version,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AvatarSrc:    user.AvatarSrc,
		BirthDate:    timestamppb.New(user.BirthDate.Time),
	}
}

func ProtoToUser(userProto *commonProto.UserResponse) (*models.User, error) {
	if err := userProto.BirthDate.CheckValid(); err != nil {
		return nil, err
	}

	return &models.User{
		ID:        userProto.Id,
		Version:   userProto.Version,
		Username:  userProto.Username,
		Email:     userProto.Email,
		Password:  userProto.PasswordHash,
		FirstName: userProto.FirstName,
		LastName:  userProto.LastName,
		AvatarSrc: userProto.AvatarSrc,
		BirthDate: models.Date{Time: userProto.BirthDate.AsTime()},
	}, nil
}
