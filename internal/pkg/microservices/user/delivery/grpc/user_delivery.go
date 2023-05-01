package grpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	proto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/user/proto/generated"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type userGRPC struct {
	userServices user.Usecase
	logger       logger.Logger

	proto.UnimplementedUserServer
}

func NewUserGRPC(userServices user.Usecase, l logger.Logger) *userGRPC {
	return &userGRPC{
		userServices: userServices,
		logger:       l,
	}
}

func (u *userGRPC) GetByID(ctx context.Context, msg *proto.Id) (*proto.UserResponse, error) {
	user, err := u.userServices.GetByID(ctx, msg.Id)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return userToProto(user), nil
}

func (u *userGRPC) UpdateInfo(ctx context.Context, msg *proto.UpdateInfoMsg) (*proto.UpdateInfoResponse, error) {
	time, err := time.Parse("2006-01-02", msg.BirthDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse date")
	}

	userInfo := &models.User{
		ID:        msg.Id,
		Email:     msg.Email,
		FirstName: msg.FirstName,
		LastName:  msg.LastName,
		Sex:       models.Sex(msg.Sex),
		BirthDate: models.Date{Time: time},
	}

	if err := u.userServices.UpdateInfo(ctx, userInfo); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.UpdateInfoResponse{}, nil
}

type uploadAvatarUsecaseInput struct {
	userID        uint32
	file          io.ReadSeeker
	fileExtension string
}

func (u *userGRPC) UploadAvatar(instream proto.User_UploadAvatarServer) error {
	usecaseInput := uploadAvatarUsecaseInput{}
	buffer := bytes.Buffer{}

	streamGotExtra := false
	for {
		req, err := instream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		switch data := req.Data.(type) {
		case *proto.UploadAvatarMsg_Extra:
			if streamGotExtra {
				return status.Error(codes.InvalidArgument, "extra message transfered double")
			}
			usecaseInput.userID = data.Extra.UserId
			usecaseInput.fileExtension = data.Extra.FileExtension
			streamGotExtra = true
		case *proto.UploadAvatarMsg_FileChunk:
			if _, err := buffer.Write(data.FileChunk); err != nil {
				return status.Error(codes.InvalidArgument, "got invalid file chunk")
			}
		}
	}

	usecaseInput.file = bytes.NewReader(buffer.Bytes())
	if err := u.userServices.UploadAvatar(instream.Context(),
		usecaseInput.userID,
		usecaseInput.file,
		usecaseInput.fileExtension); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return status.Error(codes.NotFound, err.Error())
		}

		return status.Error(codes.Internal, err.Error())
	}

	if err := instream.SendAndClose(&proto.UploadAvatarResponse{}); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (u *userGRPC) GetByPlaylist(ctx context.Context, msg *proto.GetByPlaylistMsg) (*proto.GetByPlaylistResponse, error) {
	users, err := u.userServices.GetByPlaylist(ctx, msg.PlaylistId)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	usersProto := make([]*proto.UserResponse, 0, len(users))
	for _, u := range users {
		usersProto = append(usersProto, userToProto(&u))
	}

	return &proto.GetByPlaylistResponse{Users: usersProto}, nil
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
