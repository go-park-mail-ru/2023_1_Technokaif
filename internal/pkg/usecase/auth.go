package usecase

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

const (
	salt = "@k8#&o1h18-9fwa_hg"

	secret   = "yarik_tri_is_god_of_russia"
	tokenTTL = 24 * time.Hour
)

type AuthUsecase struct {
	repo   repository.Auth
	logger logger.Logger
}

type jwtClaims struct {
	UserId 		uint `json:"id"`
	UserVersion uint `json:"user_version"`
	jwt.RegisteredClaims
}

func NewAuthUsecase(ra repository.Auth, l logger.Logger) *AuthUsecase {
	return &AuthUsecase{repo: ra, logger: l}
}

func (a *AuthUsecase) CreateUser(u models.User) (int, error) {
	u.Password = getPasswordHash(u.Password)
	return a.repo.CreateUser(u)
}

func (a *AuthUsecase) GetUserByCreds(username, password string) (*models.User, error) {
	passwordHash := getPasswordHash(password)
	user, err := a.repo.GetUserByCreds(username, passwordHash)
	if err != nil {
		return &models.User{}, err // TODO it can be repos error too
	}

	return user, nil
}

func (a *AuthUsecase) GetUserByAuthData(userID, userVersion uint) (*models.User, error) {
	return a.repo.GetUserByAuthData(userID, userVersion)
}

func (a *AuthUsecase) GenerateAccessToken(userID uint, userVersion uint) (string, error) {
	claims := &jwtClaims{
		userID,
		userVersion,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *AuthUsecase) CheckAccessToken(acessToken string) (uint, uint, error) {
	token, err := jwt.ParseWithClaims(acessToken, &jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(secret), nil
		})
	if err != nil {
		return 0, 0, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return 0, 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.UserVersion, nil
}

func (a *AuthUsecase) ChangeUserVersion(userID uint) (error) {
	err := a.repo.ChangeUserVersion(userID)
	if err != nil {
		return errors.New("failed to update user version")
	}

	return nil
}

func getPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashWithSalt := hash.Sum([]byte(salt))

	return fmt.Sprintf("%x", hashWithSalt)
}
