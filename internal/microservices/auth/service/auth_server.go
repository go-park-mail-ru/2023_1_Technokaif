package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/microservices/auth/proto"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
)

type authService struct {
	// authRepo auth.Repository
	userRepo user.Repository
	proto.UnsafeAuthorizationServer  // ?????
}

func NewAuthService(userRepo user.Repository) proto.AuthorizationServer {
	return &authService{userRepo: userRepo}
}

func (a authService) SignUpUser(ctx context.Context, msg *proto.SignUpMsg) (*proto.SignUpResponse, error) {
	salt := generateRandomSalt()

	time, err := time.Parse("2006-01-02", msg.BirthDate)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: msg.Username,
		Email: msg.Email,
		FirstName: msg.FirstName,
		LastName: msg.LastName,
		Sex: models.Sex(msg.Sex),
	}
	user.BirthDate.Time = time

	user.Salt = hex.EncodeToString(salt)

	user.Password = hashPassword(msg.Password, salt)

	userId, err := a.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &proto.SignUpResponse{Id: userId}, err
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