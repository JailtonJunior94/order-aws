package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type (
	StorageClient interface {
		GetObject(ctx context.Context, key string) ([]byte, error)
		PutObject(ctx context.Context, key string, body io.Reader) error
		SignedURL(ctx context.Context, key string, expireInSeconds int64) (string, error)
	}

	storageClient struct {
		bucket        string
		s3Client      *s3.Client
		presignClient *s3.PresignClient
	}
)

func NewStorageClient(ctx context.Context, sdkConfig aws.Config, bucketName string) (*storageClient, error) {
	if len(bucketName) == 0 {
		return nil, errors.New("s3: bucket name cannot be empty")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "http://s3.localhost.localstack.cloud:4566",
					SigningRegion: "us-east-1",
				}, nil
			},
		)),
	)

	if err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	return &storageClient{bucket: bucketName, s3Client: s3Client, presignClient: presignClient}, nil
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

func (s *storageClient) GetObject(ctx context.Context, key string) ([]byte, error) {
	resp, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("storage: failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}
	return body, nil
}

func (s *storageClient) SignedURL(ctx context.Context, key string, expireInSeconds int64) (string, error) {
	input := &s3.PutObjectInput{
		Key:         aws.String(key),
		Bucket:      aws.String(s.bucket),
		ContentType: aws.String("application/json"),
	}
	request, err := s.presignClient.PresignPutObject(ctx, input, func(p *s3.PresignOptions) {
		p.Expires = time.Duration(expireInSeconds) * time.Second
	})
	if err != nil {
		return "", fmt.Errorf("storage: %v", err)
	}
	return request.URL, nil
}
