package usecase

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/go-park-mail-ru/2023_1_Technokaif/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/repository"
)

const (
	salt = "@k8#&o1h18-9fwa_hg"

	secret   = "yarik_tri"
	tokenTTL = 24 * time.Hour
)

type AuthUsecase struct {
	repo repository.Auth
}

type jwtClaims struct {
	UserId uint `json:"id"`
	jwt.RegisteredClaims
}

func NewAuthUsecase(ra repository.Auth) *AuthUsecase {
	return &AuthUsecase{repo: ra}
}

func (a *AuthUsecase) CreateUser(u models.User) (int, error) {
	u.Password = getPasswordHash(u.Password)
	return a.repo.CreateUser(u)
}

func (a *AuthUsecase) GetUserID(username, password string) (uint, error) {
	passwordHash := getPasswordHash(password)
	user, err := a.repo.GetUser(username, passwordHash)
	if err != nil {
		return 0, errors.New("user not found") // TODO it can be repos error too
	}

	fmt.Println(user.ID)

	return user.ID, nil
}

func (a *AuthUsecase) GenerateToken(userID uint) (string, error) {
	claims := &jwtClaims{
		userID,
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

func getPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Sum([]byte(salt))

	return fmt.Sprintf("%x", hash)
}
