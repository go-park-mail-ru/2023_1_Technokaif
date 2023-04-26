package usecase

import (
	"fmt"
	"io"
	"path/filepath"

	commonFile "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/file"
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

func (u *Usecase) GetByPlaylist(playlistID uint32) ([]models.User, error) {
	users, err := u.repo.GetByPlaylist(playlistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get users of playlist: %w", err)
	}
	return users, nil
}

func (u *Usecase) UpdateInfo(user *models.User) error {
	if err := u.repo.Check(user.ID); err != nil {
		return fmt.Errorf("(usecase) can't find user with id #%d: %w", user.ID, err)
	}

	if err := u.repo.UpdateInfo(user); err != nil {
		return fmt.Errorf("(usecase) can't change user in repository: %w", err)
	}

	return nil
}

var dirForUserAvatar = filepath.Join(commonFile.MediaPath(), commonFile.AvatarFolder())

func (u *Usecase) UploadAvatar(userID uint32, file io.ReadSeeker, fileExtension string) error {
	if err := u.repo.Check(userID); err != nil {
		return fmt.Errorf("(usecase) can't find user with id #%d: %w", userID, err)
	}

	// Check format
	if fileType, err := commonFile.CheckMimeType(file, "image/png", "image/jpeg"); err != nil {
		return fmt.Errorf("(usecase) file format %s: %w", fileType, &models.AvatarWrongFormatError{FileType: fileType})
	}

	filenameWithExtension, _, err := commonFile.CreateFile(file, fileExtension, dirForUserAvatar)
	if err != nil {
		return fmt.Errorf("(usecase) can't create file: %w", err)
	}

	avatarSrc := filepath.Join(commonFile.AvatarFolder(), filenameWithExtension)
	if err := u.repo.UpdateAvatarSrc(userID, avatarSrc); err != nil {
		return fmt.Errorf("(usecase) can't update avatarSrc: %w", err)
	}
	return nil
}
