package usecase

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"

	common "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common"
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
		return &models.User{}, fmt.Errorf("(usecase) can't get user by id: %w", err)
	}
	return user, nil
}

func (u *Usecase) UpdateInfo(user *models.User) error {
	if _, err := u.repo.GetByID(user.ID); err != nil {
		return err
	}

	if err := u.repo.UpdateInfo(user); err != nil {
		return fmt.Errorf("(usecase) can't change user in repository: %w", err)
	}

	return nil
}

var dirForUserAvatars = common.MediaPath() + common.AvatarFolder()

var ErrAvatarWrongFormat = errors.New("wrong avatar file fromat")

func (u *Usecase) UploadAvatarWrongFormatError() error {
	return ErrAvatarWrongFormat
}

func (u *Usecase) UploadAvatar(user *models.User, file io.ReadSeeker, fileExtension string) error {
	// Check format
	if fileType, err := common.CheckMimeType(file, "image/png", "image/jpeg"); err != nil {
		return fmt.Errorf("(usecase) file format %s: %w", fileType, ErrAvatarWrongFormat)
	}

	// Create standard filename
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("(usecase): can't write sent avatar to hasher: %w", err)
	}
	newFileName := fmt.Sprintf("%x", hasher.Sum(nil))

	filenameWithExtencion := newFileName + "." + fileExtension

	// Save path to avatar into user entry
	path := dirForUserAvatars + "/" + filenameWithExtencion

	user.AvatarSrc = common.AvatarFolder() + "/" + filenameWithExtencion

	if _, err := os.Stat(path); os.IsNotExist(err) {
		newFD, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("(usecase): can't create file to save avatar: %w", err)
		}
		defer newFD.Close()

		if _, err := io.Copy(newFD, file); err != nil {
			return fmt.Errorf("(usecase): can't write sent avatar to file: %w", err)
		}
	}

	return u.UpdateInfo(user)
}
