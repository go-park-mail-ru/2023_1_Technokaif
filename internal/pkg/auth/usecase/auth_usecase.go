package usecase

import (
	"context"
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
	authAgent auth.Agent

	authRepo auth.Repository
	userRepo user.Repository

	logger logger.Logger
}

func NewUsecase(aa auth.Agent, ar auth.Repository, ur user.Repository, l logger.Logger) *Usecase {
	return &Usecase{
		authAgent: aa,

		authRepo: ar,
		userRepo: ur,

		logger: l,
	}
}

func (u *Usecase) SignUpUser(user models.User) (uint32, error) {
	userId, err := u.authAgent.SignUpUser(context.Background(), user) // TODO request context
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't create user: %w", err)
	}
	return userId, nil
}

func (u *Usecase) GetUserByCreds(username, password string) (*models.User, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("(usecase) cannot find user: %w", err)
	}

	valid, err := u.authAgent.CheckPassword(context.Background(), password, user.Salt, user.Password)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't check password: %w", err)
	}

	if !valid {
		return nil, fmt.Errorf("(usecase) password hash doesn't match the real one: %w",
			&models.IncorrectPasswordError{UserID: user.ID})
	}

	// salt, err := hex.DecodeString(user.Salt)
	// if err != nil {
	// 	return nil, fmt.Errorf("(usecase) invalid salt: %w", err)
	// }

	// if hashPassword(password, salt) != user.Password {
	// 	return nil, fmt.Errorf("(usecase) password hash doesn't match the real one: %w",
	// 		&models.IncorrectPasswordError{UserID: user.ID})
	// }

	return user, nil
}

func (u *Usecase) GetUserByAuthData(userID, userVersion uint32) (*models.User, error) {
	user, err := u.authRepo.GetUserByAuthData(userID, userVersion)
	if err != nil {
		return nil, fmt.Errorf("(usecase) cannot find user by id and version: %w", err)
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
	if err := u.authRepo.UpdatePassword(userID, passHash, hex.EncodeToString(salt)); err != nil {
		return fmt.Errorf("(usecase) failed to update password: %w", err)
	}

	return nil
}

func hashPassword(plainPassword string, salt []byte) string {
	hashedPassword := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return hex.EncodeToString(hashedPassword)
}

func generateRandomSalt() []byte {
	salt := make([]byte, 8)
	rand.Read(salt)
	return salt
}
