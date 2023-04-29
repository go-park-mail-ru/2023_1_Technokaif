package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	proto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/proto/generated"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type authGRPC struct {
	authServices auth.Usecase
	logger       logger.Logger

	proto.UnimplementedAuthorizationServer
}

func NewAuthGRPC(authServices auth.Usecase, l logger.Logger) proto.AuthorizationServer {
	return &authGRPC{
		authServices: authServices,
		logger:       l,
	}
}

func (a *authGRPC) SignUpUser(ctx context.Context, msg *proto.SignUpMsg) (*proto.SignUpResponse, error) {

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
		Password:  msg.Password,
	}
	user.BirthDate.Time = time

	userId, err := a.authServices.SignUpUser(ctx, user)

	if err != nil {
		var errUserAlreadyExists *models.UserAlreadyExistsError
		if errors.As(err, &errUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.SignUpResponse{UserID: userId}, nil
}

func (a *authGRPC) GetUserByCreds(ctx context.Context, msg *proto.Creds) (*proto.UserResponse, error) {
	user, err := a.authServices.GetUserByCreds(ctx, msg.Username, msg.Password)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		var errIncorrectPassword *models.IncorrectPasswordError
		if errors.As(err, &errIncorrectPassword) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return userToProto(user), nil
}

func (a *authGRPC) GetUserByAuthData(ctx context.Context, msg *proto.AuthData) (*proto.UserResponse, error) {
	user, err := a.authServices.GetUserByAuthData(ctx, msg.Id, msg.Version)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return userToProto(user), nil
}

func (a *authGRPC) IncreaseUserVersion(ctx context.Context, msg *proto.IncreaseUserVersionMsg) (*proto.Void, error) {
	if err := a.authServices.IncreaseUserVersion(ctx, msg.UserId); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (a *authGRPC) ChangePassword(ctx context.Context, msg *proto.ChangePassMsg) (*proto.Void, error) {
	if err := a.authServices.ChangePassword(ctx, msg.UserId, msg.PlainPassword); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
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
