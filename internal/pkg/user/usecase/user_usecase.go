package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements user.Usecase
type Usecase struct {
	repo   user.Repository
	logger logger.Logger
}

func NewUsecase(r user.Repository, l logger.Logger) *Usecase {
	return &Usecase{repo: r, logger: l}
}

func (u *Usecase) GetByID(userID uint32) (*models.User, error) {
	user, err := u.repo.GetByID(userID)
	if err != nil {
		return &models.User{}, fmt.Errorf("(usecase) can't get user by id : %w", err)
	}
	return user, nil
}

