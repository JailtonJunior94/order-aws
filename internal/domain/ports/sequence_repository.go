package ports

import (
	"context"

	"github.com/jailtonjunior94/order-aws/internal/domain/entities"
)

type SequenceRepository interface {
	NextValue(ctx context.Context, sequence *entities.Sequence) (int64, error)
}
