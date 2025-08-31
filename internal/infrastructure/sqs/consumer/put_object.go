package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jailtonjunior94/order-aws/internal/application/usecase"
	"github.com/jailtonjunior94/order-aws/pkg/messaging"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type PutObjectHandler struct {
	createOrderUseCase usecase.CreateOrderUseCase
}

func NewPutObjectHandler(createOrderUseCase usecase.CreateOrderUseCase) *PutObjectHandler {
	return &PutObjectHandler{
		createOrderUseCase: createOrderUseCase,
	}
}

func (h *PutObjectHandler) Handle(ctx context.Context, message types.Message) error {
	var s3EventNotification *messaging.S3EventNotifications
	if err := json.Unmarshal([]byte(*message.Body), &s3EventNotification); err != nil {
		return fmt.Errorf("failed to unmarshal S3 event notification: %v", err)
	}

	for _, record := range s3EventNotification.Records {
		if err := h.createOrderUseCase.Execute(ctx, record.S3.Object.Key); err != nil {
			return fmt.Errorf("failed to process S3 event record: %v", err)
		}
	}
	return nil
}
