package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/microservice/grpc/proto"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type AuthService struct {
	authRepo auth.Repository
	userRepo user.Repository
	logger   logger.Logger

	proto.UnimplementedAuthorizationServer
}

func NewAuthService(userRepo user.Repository, authRepo auth.Repository, l logger.Logger) proto.AuthorizationServer {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
		logger:   l,
	}
}

func (a *AuthService) SignUpUser(ctx context.Context, msg *proto.SignUpMsg) (*proto.SignUpResponse, error) {
	salt := generateRandomSalt()

	time, err := time.Parse("2006-01-02", msg.BirthDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse date")
	}

	user := models.User{
		Username:  msg.Username,
		Email:     msg.Email,
		FirstName: msg.FirstName,
		LastName:  msg.LastName,
		Sex:       models.Sex(msg.Sex),
	}
	user.BirthDate.Time = time

	user.Salt = hex.EncodeToString(salt)

	user.Password = hashPassword(msg.Password, salt)

	userId, err := a.userRepo.CreateUser(user)

	if err != nil {
		var errUserAlreadyExists *models.UserAlreadyExistsError
		if errors.As(err, &errUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &proto.SignUpResponse{UserID: userId}, nil
}

func (a *AuthService) GetUserByCreds(ctx context.Context, msg *proto.Creds) (*proto.UserResponse, error) {
	user, err := a.userRepo.GetUserByUsername(msg.Username)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to get user")
	}

	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to decode salt")
	}

	hashedPassword := hashPassword(msg.Password, salt)
	if hashedPassword != user.Password {
		return nil, status.Error(codes.PermissionDenied, "incorrect password")
	}

	return userToProto(user), nil
}

func (a *AuthService) GetUserByAuthData(ctx context.Context, msg *proto.AuthData) (*proto.UserResponse, error) {
	user, err := a.authRepo.GetUserByAuthData(ctx, msg.Id, msg.Version)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return userToProto(user), nil
}

func (a *AuthService) IncreaseUserVersion(ctx context.Context, msg *proto.IncreaseUserVersionMsg) (*proto.Void, error) {
	if err := a.authRepo.IncreaseUserVersion(ctx, msg.UserId); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to increase version")
	}

	return nil, nil
}

func (a *AuthService) ChangePassword(ctx context.Context, msg *proto.ChangePassMsg) (*proto.Void, error) {
	salt := generateRandomSalt()
	passHash := hashPassword(msg.PlainPassword, salt)
	if err := a.authRepo.UpdatePassword(ctx, msg.UserId, passHash, hex.EncodeToString(salt)); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to change pass")
	}

	return nil, nil
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

func userToProto(user *models.User) *proto.UserResponse {
	return &proto.UserResponse{
		Id:           user.ID,
		Version:      user.Version,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Sex:          string(user.Sex),
		AvatarSrc:    user.AvatarSrc,
		BirthDate:    user.BirthDate.Format("2006-01-02"),
	}
}
