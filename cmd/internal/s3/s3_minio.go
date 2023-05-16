package s3

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MakeS3MinioClient(endpoint, accessKey, secret string) (*minio.Client, error) {
	if endpoint == "" || accessKey == "" || secret == "" {
		return nil, fmt.Errorf("invalid config")
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secret, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
