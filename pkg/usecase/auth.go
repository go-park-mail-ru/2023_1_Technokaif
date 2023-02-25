package usecase

import (
	"crypto/sha256"
	"fmt"
	"time"
	"errors"

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
		return 0, err // TODO it can be repos error too
	}

	fmt.Println(user.ID)

	return user.ID, nil
}

func (a *AuthUsecase) GenerateAccessToken(userID uint) (string, error) {
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

func (a *AuthUsecase) CheckAccessToken(acessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(acessToken, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func getPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashWithSalt := hash.Sum([]byte(salt))

	return fmt.Sprintf("%x", hashWithSalt)
}
