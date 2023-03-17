package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

var secret = os.Getenv("SECRET")

const tokenTTL = 24 * time.Hour

type AuthUsecase struct {
	repo   repository.Auth
	logger logger.Logger
}

type jwtClaims struct {
	UserId      uint32 `json:"id"`
	UserVersion uint32 `json:"user_version"`
	jwt.RegisteredClaims
}

func NewAuthUsecase(ra repository.Auth, l logger.Logger) *AuthUsecase {
	return &AuthUsecase{repo: ra, logger: l}
}

func (a *AuthUsecase) CreateUser(u models.User) (uint32, error) {
	salt := make([]byte, 8)
	rand.Read(salt)
	u.Salt = fmt.Sprintf("%x", salt)

	u.Password = hashPassword(u.Password, salt)

	userId, err := a.repo.CreateUser(u)
	return userId, errors.Wrap(err, "(Usecase) cannot create user")
}

func (a *AuthUsecase) LoginUser(username, password string) (string, error) {
	user, err := a.GetUserByCreds(username, password)
	if err != nil {
		return "", errors.Wrap(err, "(Usecase) cannot find user")
	}

	token, err := a.GenerateAccessToken(user.ID, user.Version)
	if err != nil {
		return "", errors.Wrap(err, "(Usecase) failed to generate token")
	}

	return token, nil
}

func (a *AuthUsecase) GetUserByCreds(username, password string) (*models.User, error) {
	user, err := a.repo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.Wrap(err, "(Usecase) cannot find user")
	}

	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		return nil, errors.Wrap(err, "(Usecase) invalid salt")
	}

	hashedPassword := hashPassword(password, salt)
	if hashedPassword != user.Password {
		return nil, errors.New("(Usecase) password hash doesn't match the real one")
	}

	return user, nil
}

func (a *AuthUsecase) GetUserByAuthData(userID, userVersion uint32) (*models.User, error) {
	user, err := a.repo.GetUserByAuthData(userID, userVersion)
	return user, errors.Wrap(err, "(Usecase) cannot find user by userId and userVersion")
}

func (a *AuthUsecase) GenerateAccessToken(userID, userVersion uint32) (string, error) {
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
		return "", errors.Wrap(err, "(Usecase) failed to sign token")
	}

	return signedToken, nil
}

func (a *AuthUsecase) CheckAccessToken(acessToken string) (uint32, uint32, error) {
	token, err := jwt.ParseWithClaims(acessToken, &jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("(Usecase) invalid signing method")
			}
			return []byte(secret), nil
		})
	if err != nil {
		return 0, 0, errors.Wrap(err, "(Usecase) invalid token")
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return 0, 0, errors.New("(Usecase) token claims are not of type *tokenClaims")
	}

	now := time.Now().UTC()
	if claims.ExpiresAt.Time.Before(now) {
		return 0, 0, errors.New("(Usecase) token is expired")
	}

	return claims.UserId, claims.UserVersion, nil
}

func (a *AuthUsecase) IncreaseUserVersion(userID uint32) error {
	if err := a.repo.IncreaseUserVersion(userID); err != nil {
		return errors.Wrap(err, "(Usecase) failed to update user version")
	}

	return nil
}

func hashPassword(plainPassword string, salt []byte) string {
	hashedPassword := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", hashedPassword)
}
