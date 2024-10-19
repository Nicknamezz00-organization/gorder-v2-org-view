package order

import (
	"fmt"

	"github.com/Nicknamezz00/gorder-v2/common/genproto/orderpb"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v80"
)

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}

func NewOrder(id, customerID, status, paymentLink string, items []*orderpb.Item) (*Order, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	if customerID == "" {
		return nil, errors.New("empty customerID")
	}
	if status == "" {
		return nil, errors.New("empty status")
	}
	if items == nil {
		return nil, errors.New("empty items")
	}
	return &Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      status,
		PaymentLink: paymentLink,
		Items:       items,
	}, nil
}

func (o *Order) ToProto() *orderpb.Order {
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       o.Items,
		PaymentLink: o.PaymentLink,
	}
}

func (o *Order) IsPaid() error {
	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
		return nil
	}
	return fmt.Errorf("order status not paid, order id = %s, status = %s", o.ID, o.Status)
}
