package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jailtonjunior94/order-aws/pkg/storage"
)

type (
	response struct {
		URL      string `json:"url"`
		FileName string `json:"file_name"`
	}
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sdkS3Config := aws.Config{Region: "us-east-1", BaseEndpoint: aws.String("http://s3.localhost.localstack.cloud:4566")}
	s3Client, err := storage.NewStorageClient(ctx, sdkS3Config, os.Getenv("BUCKET_NAME"))
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	url, err := s3Client.SignedURL(ctx, fmt.Sprintf("testando_%s.json", time.Now().Format(time.RFC3339)), 3600)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	response := response{
		URL:      url,
		FileName: fmt.Sprintf("testando_%s.json", time.Now().Format(time.RFC3339)),
	}

	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
