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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type DynoNotation map[string]types.AttributeValue

type apiServer struct {
}

func NewApiServer() *apiServer {
	return &apiServer{}
}

func (s *apiServer) Run(ctx context.Context) {
	sdkConfig := aws.Config{Region: "us-east-1", BaseEndpoint: aws.String("http://localhost:4566")}
	dbClient, err := database.NewDynamoDBClient(ctx, sdkConfig, "local-orders")
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}

	Item := DynoNotation{
		"id":          &types.AttributeValueMemberS{Value: "123"},
		"name":        &types.AttributeValueMemberS{Value: "Sample Item"},
		"description": &types.AttributeValueMemberS{Value: "This is a sample item."},
	}

	response, err := dbClient.DynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbClient.TableName),
		Item:      Item,
	})

	if err != nil {
		log.Fatalf("Failed to put item in DynamoDB: %v", err)
	}
	log.Printf("Item added successfully: %v", response)
	log.Println("DynamoDB client initialized successfully")

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
