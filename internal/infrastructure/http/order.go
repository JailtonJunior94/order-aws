package http

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jailtonjunior94/order-aws/internal/domain/entities"
	"github.com/jailtonjunior94/order-aws/pkg/messaging"
)

type OrderHandler struct {
	sqsClient messaging.SqsClient
}

func NewOrderHandler(sqsClient messaging.SqsClient) *OrderHandler {
	return &OrderHandler{sqsClient: sqsClient}
}

func (h *OrderHandler) Handle(w http.ResponseWriter, r *http.Request) {
	entity, err := entities.NewOrder([]entities.Item{{ProductID: "123", Quantity: 2, Price: 10.0}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message, err := json.Marshal(entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.sqsClient.SendMessage(r.Context(), types.Message{Body: aws.String(string(message))}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order processed successfully"))
}
