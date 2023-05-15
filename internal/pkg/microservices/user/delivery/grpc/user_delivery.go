package grpc

import (
	"bytes"
	"context"
	"errors"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

	commonProtoUtils "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/common"
	commonProto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/common/proto/generated"
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

func (u *userGRPC) GetByID(ctx context.Context, msg *proto.Id) (*commonProto.UserResponse, error) {
	user, err := u.userServices.GetByID(ctx, msg.Id)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return commonProtoUtils.UserToProto(*user), nil
}

func (u *userGRPC) UpdateInfo(ctx context.Context, msg *proto.UpdateInfoMsg) (*proto.UpdateInfoResponse, error) {
	if err := msg.BirthDate.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userInfo := &models.User{
		ID:        msg.Id,
		Email:     msg.Email,
		FirstName: msg.FirstName,
		LastName:  msg.LastName,
		BirthDate: models.Date{Time: msg.BirthDate.AsTime()},
		Sex:       models.Male,
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
	buffer := bytes.NewBuffer(nil)

	streamGotExtra := false
	for {
		req, err := instream.Recv()
		if errors.Is(err, io.EOF) {
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
				return status.Error(codes.Internal, "got invalid file chunk")
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

	usersProto := make([]*commonProto.UserResponse, 0, len(users))
	for _, u := range users {
		usersProto = append(usersProto, commonProtoUtils.UserToProto(u))
	}

	return &proto.GetByPlaylistResponse{Users: usersProto}, nil
}
