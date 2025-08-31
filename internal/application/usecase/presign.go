package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jailtonjunior94/order-aws/internal/application/dtos"
	"github.com/jailtonjunior94/order-aws/pkg/storage"
)

type (
	PresignUseCase interface {
		Execute(ctx context.Context) (*dtos.PresignOutput, error)
	}

	presignUseCase struct {
		storage storage.StorageClient
	}
)

func NewPresignUseCase(storage storage.StorageClient) *presignUseCase {
	return &presignUseCase{storage: storage}
}

func (u *presignUseCase) Execute(ctx context.Context) (*dtos.PresignOutput, error) {
	url, err := u.storage.SignedURL(ctx, fmt.Sprintf("input_%s.json", time.Now().Format("02-01-2006_15-04-05")), 3600)
	if err != nil {
		return nil, fmt.Errorf("presign: %v", err)
	}

	return &dtos.PresignOutput{
		URL:      url,
		FileName: fmt.Sprintf("input_%s.json", time.Now().Format("02-01-2006_15-04-05")),
	}, nil
}
