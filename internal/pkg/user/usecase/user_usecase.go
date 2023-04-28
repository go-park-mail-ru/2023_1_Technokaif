package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	commonFile "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/file"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
)

// Usecase implements user.Usecase
type Usecase struct {
	repo user.Repository
}

func NewUsecase(r user.Repository) *Usecase {
	return &Usecase{
		repo: r,
	}
}

func (u *Usecase) GetByID(ctx context.Context, userID uint32) (*models.User, error) {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return &models.User{}, fmt.Errorf("(usecase) can't get user by id: %w", err)
	}
	return user, nil
}

func (u *Usecase) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.User, error) {
	users, err := u.repo.GetByPlaylist(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get users of playlist: %w", err)
	}
	return users, nil
}

func (u *Usecase) UpdateInfo(ctx context.Context, user *models.User) error {
	if err := u.repo.Check(ctx, user.ID); err != nil {
		return fmt.Errorf("(usecase) can't find user with id #%d: %w", user.ID, err)
	}

	if err := u.repo.UpdateInfo(ctx, user); err != nil {
		return fmt.Errorf("(usecase) can't change user in repository: %w", err)
	}

	return nil
}

var dirForUserAvatar = filepath.Join(commonFile.MediaPath(), commonFile.AvatarFolder())

var ErrAvatarWrongFormat = errors.New("wrong avatar file fromat")

func (u *Usecase) UploadAvatar(ctx context.Context, userID uint32, file io.ReadSeeker, fileExtension string) error {
	if err := u.repo.Check(ctx, userID); err != nil {
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
	if err := u.repo.UpdateAvatarSrc(ctx, userID, avatarSrc); err != nil {
		return fmt.Errorf("(usecase) can't update avatarSrc: %w", err)
	}
	return nil
}
