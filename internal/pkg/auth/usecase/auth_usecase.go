package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

var secret = os.Getenv("SECRET")

const tokenTTL = 24 * time.Hour

// Usecase implements auth.Usecase
type Usecase struct {
	authRepo auth.Repository
	userRepo user.Repository
	logger   logger.Logger
}

type jwtClaims struct {
	UserId      uint32 `json:"id"`
	UserVersion uint32 `json:"user_version"`
	jwt.RegisteredClaims
}

func NewUsecase(ar auth.Repository, ur user.Repository, l logger.Logger) *Usecase {
	return &Usecase{
		authRepo: ar,
		userRepo: ur,
		logger:   l}
}

func (u *Usecase) SignUpUser(user models.User) (uint32, error) {
	salt := make([]byte, 8)
	rand.Read(salt)
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
		return nil, fmt.Errorf("(usecase) password hash doesn't match the real one: %w", err)
	}

	return user, nil
}

func (u *Usecase) LoginUser(username, password string) (string, error) {
	user, err := u.GetUserByCreds(username, password)
	if err != nil {
		return "", fmt.Errorf("(usecase) cannot find user: %w", err)
	}

	token, err := u.GenerateAccessToken(user.ID, user.Version)
	if err != nil {
		return "", fmt.Errorf("(usecase) failed to generate token: %w", err)
	}

	return token, nil
}

func (u *Usecase) GetUserByAuthData(userID, userVersion uint32) (*models.User, error) {
	user, err := u.authRepo.GetUserByAuthData(userID, userVersion)
	if err != nil {
		return nil, fmt.Errorf("(usecase) cannot find user by userId and userVersion: %w", err)
	}
	return user, nil
}

func (u *Usecase) GenerateAccessToken(userID, userVersion uint32) (string, error) {
	claims := &jwtClaims{
		userID,
		userVersion,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("(usecase) failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (u *Usecase) CheckAccessToken(acessToken string) (uint32, uint32, error) {
	token, err := jwt.ParseWithClaims(acessToken, &jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("(usecase) invalid signing method")
			}
			return []byte(secret), nil
		})
	if err != nil {
		return 0, 0, fmt.Errorf("(usecase) invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return 0, 0, errors.New("(usecase) token claims are not of type *tokenClaims")
	}

	now := time.Now().UTC()
	if claims.ExpiresAt.Time.Before(now) {
		return 0, 0, errors.New("(usecase) token is expired")
	}

	return claims.UserId, claims.UserVersion, nil
}

func (u *Usecase) IncreaseUserVersion(userID uint32) error {
	if err := u.authRepo.IncreaseUserVersion(userID); err != nil {
		return fmt.Errorf("(usecase) failed to update user version: %w", err)
	}

	return nil
}

func hashPassword(plainPassword string, salt []byte) string {
	hashedPassword := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", hashedPassword)
}
