package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jailtonjunior94/order-aws/pkg/database"
	"github.com/jailtonjunior94/order-aws/pkg/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type apiServer struct {
}

func NewApiServer() *apiServer {
	return &apiServer{}
}

func (s *apiServer) Run(ctx context.Context) {
	sdkS3Config := aws.Config{Region: "us-east-1", BaseEndpoint: aws.String("http://s3.localhost.localstack.cloud:4566")}
	s3Client, err := storage.NewStorageClient(ctx, sdkS3Config, "local-orders-bucket")
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	file, err := os.Open("parametros_0.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	if err := s3Client.PutObject(ctx, "parametros_0.json", file); err != nil {
		log.Fatalf("failed to put object in S3: %v", err)
	}

	url, err := s3Client.SignedURL(ctx, "parametros_0.json", 3600)
	if err != nil {
		log.Fatalf("Failed to get signed URL: %v", err)
	}
	fmt.Printf("Signed URL: %s\n", url)

	sdkDynamoDBConfig := aws.Config{
		Region:       "us-east-1",
		BaseEndpoint: aws.String("http://localhost:4566/"),
	}
	dynamoDBClient, err := database.NewDynamoDBClient(ctx, sdkDynamoDBConfig, "local-orders")
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}

	item := database.DynoNotation{
		"id":          &types.AttributeValueMemberS{Value: "123"},
		"name":        &types.AttributeValueMemberS{Value: "Sample Item"},
		"description": &types.AttributeValueMemberS{Value: "This is a sample item."},
	}

	err = dynamoDBClient.PutItem(ctx, item)
	if err != nil {
		log.Fatalf("Failed to put item in DynamoDB: %v", err)
	}

	key := database.DynoNotation{
		"id": &types.AttributeValueMemberS{Value: "123"},
	}
	retrievedItem, err := dynamoDBClient.GetItem(ctx, key)
	if err != nil {
		log.Fatalf("Failed to get item from DynamoDB: %v", err)
	}
	fmt.Printf("Retrieved item: %v\n", retrievedItem)

	router := chi.NewRouter()
	router.Use(
		middleware.RealIP,
		middleware.RequestID,
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.AllowContentType("application/json", "application/x-www-form-urlencoded"),
	)

	/* Graceful shutdown */
	server := http.Server{
		ReadTimeout:       time.Duration(10) * time.Second,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "8001"))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	s.gracefulShutdown(&server)
}

func (s *apiServer) gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server: %v", err)
	}
	defer cancel()
}
