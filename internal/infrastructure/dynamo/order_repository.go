package dynamo

import (
	"context"
	"fmt"

	"github.com/jailtonjunior94/order-aws/internal/domain/entities"
	"github.com/jailtonjunior94/order-aws/pkg/database"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type orderRepository struct {
	dynamo database.DynamoDBClient
}

func NewOrderRepository(dynamo database.DynamoDBClient) *orderRepository {
	return &orderRepository{
		dynamo: dynamo,
	}
}

func (r *orderRepository) Save(ctx context.Context, order *entities.Order) error {
	items := make([]types.AttributeValue, len(order.Items))
	for i, item := range order.Items {
		items[i] = &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"product_id": &types.AttributeValueMemberS{Value: item.ProductID},
				"quantity":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Quantity)},
				"price":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", item.Price)},
			},
		}
	}

	item := database.DynoNotation{
		"id":    &types.AttributeValueMemberS{Value: order.ID},
		"code":  &types.AttributeValueMemberS{Value: order.Code},
		"items": &types.AttributeValueMemberL{Value: items},
	}

	if err := r.dynamo.PutItem(ctx, item); err != nil {
		return fmt.Errorf("order_repository: %v", err)
	}

	return nil
}
