package consumer

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jailtonjunior94/order-aws/configs"
	"github.com/jailtonjunior94/order-aws/internal/application/usecase"
	"github.com/jailtonjunior94/order-aws/internal/infrastructure/dynamo"
	"github.com/jailtonjunior94/order-aws/internal/infrastructure/sqs/consumer"
	"github.com/jailtonjunior94/order-aws/pkg/database"
	"github.com/jailtonjunior94/order-aws/pkg/messaging"
	"github.com/jailtonjunior94/order-aws/pkg/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func Run(ctx context.Context, config *configs.Config) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("received signal: %s. shutting down gracefully...", sig)
		cancel()
	}()

	awsConfig := aws.Config{Region: config.AWSConfig.Region, BaseEndpoint: aws.String(config.AWSConfig.Endpoint)}
	sqsClient, err := messaging.NewSqsClient(ctx, awsConfig, config.SQSConfig.QueueName)
	if err != nil {
		log.Fatalf("failed to create SQS client: %v", err)
	}

	s3ClientConfig := aws.Config{Region: config.AWSConfig.Region, BaseEndpoint: aws.String(config.AWSConfig.S3Endpoint)}
	s3Client, err := storage.NewStorageClient(ctx, s3ClientConfig, config.S3Config.BucketName)
	if err != nil {
		log.Fatalf("failed to create S3 client: %v", err)
	}

	dynamoClient, err := database.NewDynamoDBClient(ctx, awsConfig, config.DynamoDBConfig.TableName)
	if err != nil {
		log.Fatalf("failed to create DynamoDB client: %v", err)
	}

	orderRepository := dynamo.NewOrderRepository(dynamoClient)
	createOrderUseCase := usecase.NewCreateOrderUseCase(s3Client, orderRepository)
	createOrderHandler := consumer.NewPutObjectHandler(createOrderUseCase)

	go func() {
		for {
			if err := sqsClient.ReceiveMessages(
				ctx,
				config.SQSConfig.MaxNumberOfMessages,
				config.SQSConfig.WaitTimeSeconds,
				config.SQSConfig.VisibilityTimeout,
				createOrderHandler.Handle,
			); err != nil {
				log.Printf("failed to receive messages: %v", err)
			}
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
		}
	}()

	<-ctx.Done()
	log.Println("consumer has been shut down.")
}
