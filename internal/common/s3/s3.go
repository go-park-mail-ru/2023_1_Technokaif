package s3

/*import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type Buckets struct {
	avatarBucket   string
	playlistBucket string
}

type S3Saver struct {
	buckets Buckets
	cl      *minio.Client
}

func NewS3Saver(buckets Buckets, client *minio.Client) *S3Saver {
	return &S3Saver{
		buckets: buckets,
		cl:      client,
	}
}

func (s *S3Saver) SaveAvatar(ctx context.Context, avatar io.Reader, objectName string, size int64) error {
	uploadInfo, err := s.cl.PutObject(context.Background(), s.buckets.avatarBucket, objectName,
						file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
} */
