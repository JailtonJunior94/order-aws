package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jailtonjunior94/order-aws/internal/application/usecase"
	"github.com/jailtonjunior94/order-aws/internal/domain/entities"
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
	var eventNotification *messaging.EventNotifications
	if err := json.Unmarshal([]byte(*message.Body), &eventNotification); err == nil && len(eventNotification.Records) > 0 {
		for _, record := range eventNotification.Records {
			if err := h.createOrderUseCase.Execute(ctx, record.S3.Object.Key); err != nil {
				return fmt.Errorf("failed to process S3 event record: %v", err)
			}
			fmt.Printf("Processed S3 event record: %+v\n", record)
		}
		return nil
	}

	var eventRecord *entities.Order
	if err := json.Unmarshal([]byte(*message.Body), &eventRecord); err != nil {
		return fmt.Errorf("failed to unmarshal order event: %v", err)
	}

	fmt.Printf("Order processed: %+v\n", eventRecord)
	return nil
}
