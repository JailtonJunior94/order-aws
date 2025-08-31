package consumer

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jailtonjunior94/order-aws/configs"
	"github.com/jailtonjunior94/order-aws/pkg/messaging"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func Run(config *configs.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %s. Shutting down gracefully...", sig)
		cancel()
	}()

	awsConfig := aws.Config{Region: "us-east-1", BaseEndpoint: aws.String("http://localhost:4566/")}
	sqsClient, err := messaging.NewSqsClient(ctx, awsConfig, config.SQSConfig.QueueName)
	if err != nil {
		log.Fatalf("Failed to create SQS client: %v", err)
	}

	go func() {
		for {
			if err := sqsClient.ReceiveMessages(ctx, 10, handlerMessage); err != nil {
				log.Fatalf("Failed to receive messages: %v", err)
			}
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	<-ctx.Done()
	log.Println("Consumer has been shut down.")
}

func handlerMessage(ctx context.Context, body []byte) error {
	log.Println("Message received: ", string(body))
	if string(body) == "error" {
		return errors.New("deu ruim")
	}
	return nil
}
