package ports

import (
	"context"

	"github.com/jailtonjunior94/order-aws/internal/domain/entities"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entities.Order) error
}
