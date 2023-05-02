package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	commonProtoUtils "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/common"
	proto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/user/proto/generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserAgent struct {
	client proto.UserClient
}

func NewUserAgent(client proto.UserClient) *UserAgent {
	return &UserAgent{
		client: client,
	}
}

func (u *UserAgent) GetByID(ctx context.Context, userID uint32) (*models.User, error) {
	resp, err := u.client.GetByID(ctx, &proto.Id{Id: userID})
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

	user, err := commonProtoUtils.ProtoToUser(resp)
	if err != nil {
		return nil, fmt.Errorf("(usecase) convert from proto to user: %w", err)
	}

	return user, nil
}

func (u *UserAgent) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.User, error) {
	resp, err := u.client.GetByPlaylist(ctx, &proto.GetByPlaylistMsg{PlaylistId: playlistID})
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

	users := make([]models.User, 0, len(resp.Users))
	for _, protoUser := range resp.Users {
		user, err := commonProtoUtils.ProtoToUser(protoUser)
		if err != nil {
			return nil, fmt.Errorf("(usecase) convert from proto to user: %w", err)
		}

		users = append(users, *user)
	}

	return users, nil
}

func (u *UserAgent) UpdateInfo(ctx context.Context, user *models.User) error {
	_, err := u.client.UpdateInfo(ctx, userToProtoUserInfo(user))
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return fmt.Errorf("%w: %v", &models.NoSuchUserError{}, err)
			case codes.Internal, codes.InvalidArgument:
				return err
			}
		}
		return err
	}

	return nil
}

func (u *UserAgent) UploadAvatar(ctx context.Context, userID uint32, file io.ReadSeeker, fileExtension string) error {
	stream, err := u.client.UploadAvatar(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&proto.UploadAvatarMsg{
		Data: &proto.UploadAvatarMsg_Extra{
			Extra: &proto.UploadAvatarExtra{
				UserId:        userID,
				FileExtension: fileExtension,
			},
		},
	})
	if err != nil {
		return nil
	}

	fileChunk := make([]byte, 1024)

	for {
		bytesRead, err := file.Read(fileChunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("(usecase) failed to send chunk")
		}

		err = stream.Send(&proto.UploadAvatarMsg{
			Data: &proto.UploadAvatarMsg_FileChunk{
				FileChunk: fileChunk[:bytesRead],
			},
		})
		if err != nil {
			return fmt.Errorf("(usecase) failed to send chunk")
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return fmt.Errorf("%w: %v", &models.NoSuchUserError{}, err)
			case codes.Internal, codes.InvalidArgument:
				return err
			}
		}
		return err
	}

	return nil
}

func userToProtoUserInfo(user *models.User) *proto.UpdateInfoMsg {
	return &proto.UpdateInfoMsg{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Sex:       string(user.Sex),
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}
}
