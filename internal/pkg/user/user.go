package user

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

type UserUsecase interface {
	GetByID(userID uint32) (models.UserTransfer, error)
}

type UserRepository interface {
	GetByID(userID uint32) (models.User, error)
}