package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"

	proto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/proto/generated"
	commonProtoUtils "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/common"
	commonProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/common/proto/generated"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type authGRPC struct {
	authServices auth.Usecase
	logger       logger.Logger

	proto.UnimplementedAuthorizationServer
}

func NewAuthGRPC(authServices auth.Usecase, l logger.Logger) *authGRPC {
	return &authGRPC{
		authServices: authServices,
		logger:       l,
	}
}

func (a *authGRPC) SignUpUser(ctx context.Context, msg *proto.SignUpMsg) (*proto.SignUpResponse, error) {
	if err := msg.BirthDate.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := models.User{
		Username:  msg.Username,
		Email:     msg.Email,
		FirstName: msg.FirstName,
		LastName:  msg.LastName,
		Password:  msg.Password,
		BirthDate: models.Date{Time: msg.BirthDate.AsTime()},
	}

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

func (a *authGRPC) GetUserByCreds(ctx context.Context, msg *proto.Creds) (*commonProto.UserResponse, error) {
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

	return commonProtoUtils.UserToProto(*user), nil
}

func (a *authGRPC) GetUserByAuthData(ctx context.Context, msg *proto.AuthData) (*commonProto.UserResponse, error) {
	user, err := a.authServices.GetUserByAuthData(ctx, msg.Id, msg.Version)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return commonProtoUtils.UserToProto(*user), nil
}

func (a *authGRPC) IncreaseUserVersion(ctx context.Context, msg *proto.IncreaseUserVersionMsg) (*proto.IncreaseUserVersionResponse, error) {
	if err := a.authServices.IncreaseUserVersion(ctx, msg.UserId); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.IncreaseUserVersionResponse{}, nil
}

func (a *authGRPC) ChangePassword(ctx context.Context, msg *proto.ChangePassMsg) (*proto.ChangePassResponse, error) {
	if err := a.authServices.ChangePassword(ctx, msg.UserId, msg.PlainPassword); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.ChangePassResponse{}, nil
}
