package usecase_test

import (
	"context"
	"testing"

	"github.com/jailtonjunior94/order-aws/internal/application/usecase"
	orderRepositoryMock "github.com/jailtonjunior94/order-aws/internal/infrastructure/dynamo/mocks"
	storageMock "github.com/jailtonjunior94/order-aws/pkg/storage/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateOrderUseCaseSuite struct {
	suite.Suite

	ctx             context.Context
	storage         *storageMock.StorageClient
	orderRepository *orderRepositoryMock.OrderRepository
}

func TestCreateOrderUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderUseCaseSuite))
}

func (s *CreateOrderUseCaseSuite) SetupTest() {
	s.storage = storageMock.NewStorageClient(s.T())
	s.orderRepository = orderRepositoryMock.NewOrderRepository(s.T())
}

func (s *CreateOrderUseCaseSuite) TestExecute() {
	type (
		args struct {
			input string
		}

		dependencies struct {
			storage         *storageMock.StorageClient
			orderRepository *orderRepositoryMock.OrderRepository
		}
	)

	jsonBytes := []byte(`{
    "id": "12345",
    "items": [
        {
            "product_id": "111",
            "quantity": 2,
            "price": 10.0   
        },
        {
            "product_id": "222",
            "quantity": 1,
            "price": 20.0
        }
    ]
}`)

	scenarios := []struct {
		name         string
		args         args
		dependencies dependencies
		expect       func(err error)
	}{
		{
			name: "deve receber evento de arquivo criado e salvar no dynamo",
			args: args{input: "meu-arquivo.json"},
			dependencies: dependencies{
				storage: func() *storageMock.StorageClient {
					s.storage.
						EXPECT().
						GetObject(s.ctx, mock.Anything).
						Return(jsonBytes, nil)
					return s.storage
				}(),
				orderRepository: func() *orderRepositoryMock.OrderRepository {
					s.orderRepository.
						EXPECT().
						Save(s.ctx, mock.Anything).
						Return(nil)
					return s.orderRepository
				}(),
			},
			expect: func(err error) {
				s.NoError(err)
			},
		},
	}

	for _, scenario := range scenarios {
		s.T().Run(scenario.name, func(t *testing.T) {
			createOrderUseCase := usecase.NewCreateOrderUseCase(
				scenario.dependencies.storage,
				scenario.dependencies.orderRepository,
			)
			err := createOrderUseCase.Execute(s.ctx, scenario.args.input)
			scenario.expect(err)
		})
	}
}
