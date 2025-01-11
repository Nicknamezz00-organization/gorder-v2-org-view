package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nicknamezz00/gorder-v2/common/entity"
	"github.com/Nicknamezz00/gorder-v2/common/tracing"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

type StripeProcessor struct {
	apiKey string
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	if apiKey == "" {
		panic("empty api key")
	}
	stripe.Key = apiKey
	return &StripeProcessor{apiKey: apiKey}
}

const (
	successURL = "http://localhost:8282/success"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *entity.Order) (string, error) {
	_, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
	defer span.End()

	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	marshalledItems, _ := json.Marshal(order.Items)
	metadata := map[string]string{
		"orderID":     order.ID,
		"customerID":  order.CustomerID,
		"status":      order.Status,
		"items":       string(marshalledItems),
		"paymentLink": order.PaymentLink,
	}
	params := &stripe.CheckoutSessionParams{
		Metadata:   metadata,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s?customerID=%s&orderID=%s", successURL, order.CustomerID, order.ID)),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
