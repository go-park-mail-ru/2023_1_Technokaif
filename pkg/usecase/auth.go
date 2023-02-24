package usecase

import (
	"crypto/sha256"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/repository"
)

const (
	salt   = "aknfio1h189fwahg"
	secret = "yarik_tri"
)

type AuthUsecase struct {
	repo repository.Auth
}

func NewAuthUsecase(ra repository.Auth) *AuthUsecase {
	return &AuthUsecase{repo: ra}
}

func (a *AuthUsecase) CreateUser(u models.User) (int, error) {
	u.Password = passwordHash(u.Password)
	return a.repo.CreateUser(u)
}

func (a *AuthUsecase) GenerateToken(username, password string) (string, error) {
	return "", nil
}

func passwordHash(passwd string) string {
	hash := sha256.New()
	hash.Write([]byte(passwd))
	hashWithSalt := sha256.Sum256([]byte(salt))
	return fmt.Sprintf("%x", hashWithSalt)
}
