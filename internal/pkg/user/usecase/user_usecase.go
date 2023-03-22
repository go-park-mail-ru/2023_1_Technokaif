package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/google/uuid"
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

func (u *Usecase) ChangeInfo(user *models.User) error {
	if err := u.repo.UpdateInfo(user); err != nil {
		return fmt.Errorf("(usecase) can't change user in repository: %w", err)
	}

	return nil
}

const dirForUserAvatars = "./img/user_avatars"

func (u *Usecase) UploadAvatar(user *models.User, fileHeader *multipart.FileHeader) error {
	fd, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("(usecase): can't open file of sent avatar: %w", err)
	}
	defer fd.Close()

	// Create standard filename
	newFileName := "user" + strconv.Itoa(int(user.ID)) + "_" + uuid.NewString()
	filenameParts := strings.Split(fileHeader.Filename, ".")
	fileExtension := filenameParts[len(filenameParts)-1]

	filename := newFileName + "." + fileExtension

	// Save path to avatar into user entry
	path := dirForUserAvatars + "/" + filename
	user.AvatarSrc = path

	newFD, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("(usecase): can't create file to save avatar: %w", err)
	}
	defer newFD.Close()

	if _, err := io.Copy(newFD, fd); err != nil {
		return fmt.Errorf("(usecase): can't copy sent avatar: %w", err)
	}

	return u.ChangeInfo(user)
}
