package main

import (
	"context"
	"log"
	"os"

	"github.com/jailtonjunior94/order-aws/internal/application/usecase"
	"github.com/jailtonjunior94/order-aws/internal/infrastructure/http"
	"github.com/jailtonjunior94/order-aws/pkg/storage"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {
	sdkS3Config := aws.Config{Region: os.Getenv("REGION"), BaseEndpoint: aws.String(os.Getenv("BUCKET_ENDPOINT"))}
	storage, err := storage.NewStorageClient(context.Background(), sdkS3Config, os.Getenv("BUCKET_NAME"))
	if err != nil {
		log.Fatalf("failed to create S3 client: %v", err)
	}

	presignUseCase := usecase.NewPresignUseCase(storage)
	presignHandler := http.NewPresignHandler(presignUseCase)

	lambda.Start(presignHandler.Handle)
}
