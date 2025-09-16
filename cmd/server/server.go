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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jailtonjunior94/order-aws/configs"
	handlers "github.com/jailtonjunior94/order-aws/internal/infrastructure/http"
	"github.com/jailtonjunior94/order-aws/pkg/messaging"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(ctx context.Context, config *configs.Config) {
	router := chi.NewRouter()
	router.Use(
		middleware.RealIP,
		middleware.RequestID,
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.AllowContentType("application/json", "application/x-www-form-urlencoded"),
	)

	awsConfig := aws.Config{Region: config.AWSConfig.Region, BaseEndpoint: aws.String(config.AWSConfig.Endpoint)}
	sqsClient, err := messaging.NewSqsClient(ctx, awsConfig, config.SQSConfig.QueueName)
	if err != nil {
		log.Fatalf("failed to create SQS client: %v", err)
	}

	orderHandler := handlers.NewOrderHandler(sqsClient)

	router.Post("/orders", orderHandler.Handle)

	/* Graceful shutdown */
	server := http.Server{
		ReadTimeout:       time.Duration(10) * time.Second,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "8001"))
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	gracefulShutdown(&server)
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server: %v", err)
	}
	defer cancel()
}
