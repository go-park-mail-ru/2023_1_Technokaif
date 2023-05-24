package s3

import (
	"context"
	"io"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

type S3AvatarSaver struct {
	avatarBucket string
	avatarFolder string
	cl           *minio.Client
}

func NewS3AvatarSaver(avatarBucket, avatarFolder string, client *minio.Client) *S3AvatarSaver {
	return &S3AvatarSaver{
		avatarBucket: avatarBucket,
		avatarFolder: avatarFolder,
		cl:           client,
	}
}

func (s *S3AvatarSaver) Save(ctx context.Context, avatar io.Reader, fileName string, size int64) error {

	objectPath := filepath.Join(s.avatarFolder, fileName)
	_, err := s.cl.PutObject(context.Background(), s.avatarBucket, objectPath,
		avatar, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})

	if err != nil {
		return err
	}

	return nil
}
