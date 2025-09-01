package http

import (
	"context"
	"encoding/json"

	"github.com/jailtonjunior94/order-aws/internal/application/usecase"

	"github.com/aws/aws-lambda-go/events"
)

type PresignHandler struct {
	presignUseCase usecase.PresignUseCase
}

func NewPresignHandler(presignUseCase usecase.PresignUseCase) *PresignHandler {
	return &PresignHandler{
		presignUseCase: presignUseCase,
	}
}

func (h *PresignHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	output, err := h.presignUseCase.Execute(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	body, err := json.Marshal(output)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
