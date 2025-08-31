package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jailtonjunior94/order-aws/internal/domain/entities"
	"github.com/jailtonjunior94/order-aws/internal/domain/ports"
	"github.com/jailtonjunior94/order-aws/pkg/storage"
)

type (
	CreateOrderUseCase interface {
		Execute(ctx context.Context, objectKey string) error
	}

	createOrderUseCase struct {
		storage         storage.StorageClient
		orderRepository ports.OrderRepository
	}
)

func NewCreateOrderUseCase(
	storage storage.StorageClient,
	orderRepository ports.OrderRepository,
) *createOrderUseCase {
	return &createOrderUseCase{
		storage:         storage,
		orderRepository: orderRepository,
	}
}

func (u *createOrderUseCase) Execute(ctx context.Context, objectKey string) error {
	body, err := u.storage.GetObject(ctx, objectKey)
	if err != nil {
		return fmt.Errorf("create_order_use_case: failed to get object from storage: %v", err)
	}

	order := &entities.Order{}
	if err := json.Unmarshal(body, order); err != nil {
		return fmt.Errorf("create_order_use_case: failed to unmarshal order: %v", err)
	}

	newOrder, err := entities.NewOrder(order.Items)
	if err != nil {
		return fmt.Errorf("create_order_use_case: failed to create new order: %v", err)
	}

	if err := u.orderRepository.Save(ctx, newOrder); err != nil {
		return fmt.Errorf("create_order_use_case: failed to save order: %v", err)
	}

	return nil
}
