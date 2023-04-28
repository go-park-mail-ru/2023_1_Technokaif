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
)

// Usecase implements auth.Usecase
type Usecase struct {
	authRepo auth.Repository
	userRepo user.Repository
}

func NewUsecase(ar auth.Repository, ur user.Repository) *Usecase {
	return &Usecase{
		authRepo: ar,
		userRepo: ur,
	}
}

func (u *Usecase) SignUpUser(ctx context.Context, user models.User) (uint32, error) {
	salt := generateRandomSalt()
	user.Salt = hex.EncodeToString(salt)

	user.Password = hashPassword(user.Password, salt)

	userId, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("(usecase) cannot create user: %w", err)
	}
	return userId, nil
}

func (u *Usecase) GetUserByCreds(ctx context.Context, username, password string) (*models.User, error) {
	user, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("(usecase) cannot find user: %w", err)
	}

	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		return nil, fmt.Errorf("(usecase) invalid salt: %w", err)
	}

	hashedPassword := hashPassword(password, salt)
	if hashedPassword != user.Password {
		return nil, fmt.Errorf("(usecase) password hash doesn't match the real one: %w", &models.IncorrectPasswordError{UserID: user.ID})
	}

	return user, nil
}

func (u *Usecase) GetUserByAuthData(ctx context.Context, userID, userVersion uint32) (*models.User, error) {
	user, err := u.authRepo.GetUserByAuthData(ctx, userID, userVersion)
	if err != nil {
		return nil, fmt.Errorf("(usecase) cannot find user by userId and userVersion: %w", err)
	}
	return user, nil
}

func (u *Usecase) IncreaseUserVersion(ctx context.Context, userID uint32) error {
	if err := u.authRepo.IncreaseUserVersion(ctx, userID); err != nil {
		return fmt.Errorf("(usecase) failed to update user version: %w", err)
	}

	return nil
}

func (u *Usecase) ChangePassword(ctx context.Context, userID uint32, password string) error {
	salt := generateRandomSalt()
	passHash := hashPassword(password, salt)
	if err := u.authRepo.UpdatePassword(ctx, userID, passHash, hex.EncodeToString(salt)); err != nil {
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