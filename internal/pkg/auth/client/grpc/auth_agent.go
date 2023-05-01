package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/auth/proto/generated"
)

type AuthAgent struct {
	client proto.AuthorizationClient
}

func NewAuthAgent(c proto.AuthorizationClient) *AuthAgent {
	return &AuthAgent{
		client: c,
	}
}

func (a *AuthAgent) SignUpUser(ctx context.Context, u models.User) (uint32, error) {
	msg := &proto.SignUpMsg{
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Sex:       string(u.Sex),
		BirthDate: u.BirthDate.Format("2006-01-02"), // TODO Check format
	}

	resp, err := a.client.SignUpUser(ctx, msg)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return 0, fmt.Errorf("%w: %v", &models.UserAlreadyExistsError{}, err)
			case codes.Internal, codes.InvalidArgument:
				return 0, err
			}
		}
		return 0, err
	}

	return resp.UserID, nil
}

func (a *AuthAgent) GetUserByCreds(ctx context.Context, username, plainPassword string) (*models.User, error) {
	msg := &proto.Creds{
		Username: username,
		Password: plainPassword,
	}

	resp, err := a.client.GetUserByCreds(ctx, msg)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, fmt.Errorf("%w: %v", &models.NoSuchUserError{}, err)
			case codes.PermissionDenied:
				return nil, fmt.Errorf("%w: %v", &models.IncorrectPasswordError{}, err)
			case codes.Internal:
				return nil, err
			}
		}
		return nil, err
	}

	user, err := protoToUser(resp)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthAgent) GetUserByAuthData(ctx context.Context, userID, userVersion uint32) (*models.User, error) {
	msg := &proto.AuthData{
		Id:      userID,
		Version: userVersion,
	}

	resp, err := a.client.GetUserByAuthData(ctx, msg)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, fmt.Errorf("%w: %v", &models.NoSuchUserError{}, err)
			case codes.Internal:
				return nil, err
			}
		}
		return nil, err
	}

	user, err := protoToUser(resp)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthAgent) IncreaseUserVersion(ctx context.Context, userID uint32) error {
	_, err := a.client.IncreaseUserVersion(ctx, &proto.IncreaseUserVersionMsg{UserId: userID})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return fmt.Errorf("%w: %v", &models.NoSuchUserError{}, err)
			case codes.Internal:
				return err
			}
		}
		return err
	}

	return nil
}

func (a *AuthAgent) ChangePassword(ctx context.Context, userID uint32, password string) error {
	msg := &proto.ChangePassMsg{
		UserId:        userID,
		PlainPassword: password,
	}

	_, err := a.client.ChangePassword(ctx, msg)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return fmt.Errorf("%w: %v", &models.NoSuchUserError{}, err)
			case codes.Internal:
				return err
			}
		}
		return err
	}
	return nil
}

func protoToUser(userProto *proto.UserResponse) (*models.User, error) {
	time, err := time.Parse("2006-01-02", userProto.BirthDate)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        userProto.Id,
		Version:   userProto.Version,
		Username:  userProto.Username,
		Email:     userProto.Email,
		Password:  userProto.PasswordHash,
		FirstName: userProto.FirstName,
		LastName:  userProto.LastName,
		Sex:       models.Sex(userProto.Sex),
		AvatarSrc: userProto.AvatarSrc,
		BirthDate: models.Date{Time: time},
	}

	return user, nil
}
