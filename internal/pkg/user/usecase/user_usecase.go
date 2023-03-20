package user_usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type userUsecase struct {
	repo   user.UserRepository
	logger logger.Logger
}

func NewUserUsecase(r user.UserRepository, l logger.Logger) user.UserUsecase {
	return &userUsecase{repo: r, logger: l}
}

func (u *userUsecase) GetByID(userID uint32) (models.UserTransfer, error) {
	user, err := u.repo.GetByID(userID)
	if err != nil {
		return models.UserTransfer{}, fmt.Errorf("(usecase) can't get user by id : %w", err)
	}
	return models.UserTransfer{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Sex:       user.Sex,
		BirhDate:  user.BirhDate,
		AvatarSrc: user.AvatarSrc,
	}, nil
}
