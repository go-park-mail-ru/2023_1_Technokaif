package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements auth.Usecase
type Usecase struct {
	authRepo auth.Repository
	userRepo user.Repository
	logger   logger.Logger
}

func NewUsecase(ar auth.Repository, ur user.Repository, l logger.Logger) *Usecase {
	return &Usecase{
		authRepo: ar,
		userRepo: ur,
		logger:   l}
}

func (u *Usecase) SignUpUser(user models.User) (uint32, error) {
	salt := generateRandomSalt()
	user.Salt = fmt.Sprintf("%x", salt)

	user.Password = hashPassword(user.Password, salt)

	userId, err := u.userRepo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("(usecase) cannot create user: %w", err)
	}
	return userId, nil
}

func (u *Usecase) GetUserByCreds(username, password string) (*models.User, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("(usecase) cannot find user: %w", err)
	}

	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		return nil, fmt.Errorf("(usecase) invalid salt: %w", err)
	}

	hashedPassword := hashPassword(password, salt)
	if hashedPassword != user.Password {
		return nil, fmt.Errorf("(usecase) password hash doesn't match the real one: %w", &models.IncorrectPasswordError{UserId: user.ID})
	}

	return user, nil
}

func (u *Usecase) GetUserByAuthData(userID, userVersion uint32) (*models.User, error) {
	user, err := u.authRepo.GetUserByAuthData(userID, userVersion)
	if err != nil {
		return nil, fmt.Errorf("(usecase) cannot find user by userId and userVersion: %w", err)
	}
	return user, nil
}

func (u *Usecase) IncreaseUserVersion(userID uint32) error {
	if err := u.authRepo.IncreaseUserVersion(userID); err != nil {
		return fmt.Errorf("(usecase) failed to update user version: %w", err)
	}

	return nil
}

func (u *Usecase) ChangePassword(userID uint32, password string) error {
	salt := generateRandomSalt()
	passHash := hashPassword(password, salt)
	if err := u.authRepo.UpdatePassword(userID, passHash, fmt.Sprintf("%x", salt)); err != nil {
		return fmt.Errorf("(usecase) failed to update password: %w", err)
	}

	return nil
}

func hashPassword(plainPassword string, salt []byte) string {
	hashedPassword := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", hashedPassword)
}

func generateRandomSalt() []byte {
	salt := make([]byte, 8)
	rand.Read(salt)
	return salt
}
