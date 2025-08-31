package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/jailtonjunior94/order-aws/internal/application/usecase"
	"github.com/jailtonjunior94/order-aws/pkg/storage"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sdkS3Config := aws.Config{Region: os.Getenv("REGION"), BaseEndpoint: aws.String(os.Getenv("BUCKET_ENDPOINT"))}
	s3Client, err := storage.NewStorageClient(ctx, sdkS3Config, os.Getenv("BUCKET_NAME"))
	if err != nil {
		log.Fatalf("failed to create S3 client: %v", err)
	}

	presignUseCase := usecase.NewPresignUseCase(s3Client)
	output, err := presignUseCase.Execute(ctx)
	if err != nil {
		log.Fatalf("failed to execute presign use case: %v", err)
	}

	body, err := json.Marshal(output)
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
