package messaging

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type (
	ConsumeHandler func(ctx context.Context, message types.Message) error

	SqsClient interface {
		SendMessage(ctx context.Context, message types.Message) error
		DeleteMessage(ctx context.Context, message types.Message) error
		ReceiveMessages(ctx context.Context, maxNumberOfMessages int32, waitTimeSeconds int32, visibilityTimeout int32, handler ConsumeHandler) error
	}

	sqsClient struct {
		queue     string
		sqsClient *sqs.Client
	}
)

func NewSqsClient(ctx context.Context, sdkConfig aws.Config, queueName string) (*sqsClient, error) {
	if len(queueName) == 0 {
		return nil, errors.New("sqs: queue name cannot be empty")
	}
	return &sqsClient{queue: queueName, sqsClient: sqs.NewFromConfig(sdkConfig)}, nil
}

func (s *sqsClient) ReceiveMessages(ctx context.Context, maxNumberOfMessages int32, waitTimeSeconds int32, visibilityTimeout int32, handler ConsumeHandler) error {
	output, err := s.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queue),
		MaxNumberOfMessages: maxNumberOfMessages,
		WaitTimeSeconds:     waitTimeSeconds,
		VisibilityTimeout:   visibilityTimeout,
	})
	if err != nil {
		return fmt.Errorf("sqs: failed to receive messages: %v", err)
	}

	go func() {
		for _, message := range output.Messages {
			if err := handler(ctx, message); err != nil {
				log.Fatalf("sqs: failed to process message: %v", err)
				continue
			}

			if err := s.DeleteMessage(ctx, message); err != nil {
				log.Fatalf("sqs: failed to delete message: %v", err)
				continue
			}
		}
	}()

	return nil
}

func (s *sqsClient) SendMessage(ctx context.Context, message types.Message) error {
	_, err := s.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: message.Body,
		QueueUrl:    aws.String(s.queue),
	})
	if err != nil {
		return fmt.Errorf("sqs: failed to send message: %v", err)
	}
	return nil
}

func (s *sqsClient) DeleteMessage(ctx context.Context, message types.Message) error {
	_, err := s.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queue),
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return fmt.Errorf("sqs: failed to delete message: %v", err)
	}
	return nil
}
