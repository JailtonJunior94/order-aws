package storage

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type (
	StorageClient interface {
		PutObject(ctx context.Context) error
	}

	storageClient struct {
		bucket   string
		s3Client *s3.Client
	}
)

func NewStorageClient(ctx context.Context, sdkConfig aws.Config, bucketName string) (*storageClient, error) {
	if len(bucketName) == 0 {
		return nil, errors.New("s3: bucket name cannot be empty")
	}
	return &storageClient{bucket: bucketName, s3Client: s3.NewFromConfig(sdkConfig)}, nil
}

func (s *storageClient) PutObject(ctx context.Context, key string, body io.Reader) error {
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return err
	}
	return nil
}
