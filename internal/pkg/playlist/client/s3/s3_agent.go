package s3

import (
	"context"
	"io"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

type S3PlaylistCoverSaver struct {
	playlistCoverBucket string
	playlistCoverFolder string
	cl                  *minio.Client
}

func NewS3PlaylistCoverSaver(playlistCoverBucket, playlistCoverFolder string, client *minio.Client) *S3PlaylistCoverSaver {
	return &S3PlaylistCoverSaver{
		playlistCoverBucket: playlistCoverBucket,
		playlistCoverFolder: playlistCoverFolder,
		cl:                  client,
	}
}

func (s *S3PlaylistCoverSaver) Save(ctx context.Context, cover io.Reader, fileName string, size int64) error {

	objectPath := filepath.Join(s.playlistCoverFolder, fileName)
	_, err := s.cl.PutObject(context.Background(), s.playlistCoverBucket, objectPath,
		cover, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})

	if err != nil {
		return err
	}

	return nil
}
