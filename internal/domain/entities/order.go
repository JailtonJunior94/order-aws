package entities

import (
	"fmt"

	"github.com/JailtonJunior94/devkit-go/pkg/vos"
)

type (
	Order struct {
		ID    string `json:"id"`
		Code  string `json:"code"`
		Items []Item `json:"items"`
	}

	Item struct {
		ProductID string  `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	}
)

func NewOrder(sequence int64, items []Item) (*Order, error) {
	id, err := vos.NewULID()
	if err != nil {
		return nil, err
	}

	code := fmt.Sprintf("order#%d", sequence)
	return &Order{ID: id.String(), Code: code, Items: items}, nil
}

func (o *Order) TotalAmount() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
