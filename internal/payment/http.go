package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Nicknamezz00/gorder-v2/common/broker"
	"github.com/Nicknamezz00/gorder-v2/common/entity"
	"github.com/Nicknamezz00/gorder-v2/common/logging"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
	"go.opentelemetry.io/otel"
)

type PaymentHandler struct {
	channel *amqp.Channel
}

func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
	return &PaymentHandler{channel: ch}
}

// stripe listen --forward-to localhost:8284/api/webhook
func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.handleWebhook)
}

func (h *PaymentHandler) handleWebhook(c *gin.Context) {
	logrus.WithContext(c.Request.Context()).Info("receive webhook from stripe")
	var err error
	defer func() {
		if err != nil {
			logging.Warnf(c.Request.Context(), nil, "handleWebhook err=%v", err)
		} else {
			logging.Infof(c.Request.Context(), nil, "%s", "handleWebhook success")
		}
	}()

	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		err = errors.Wrap(err, "Error reading request body")
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}

	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("ENDPOINT_STRIPE_SECRET"))

	if err != nil {
		err = errors.Wrap(err, "error verifying webhook signature")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err = json.Unmarshal(event.Data.Raw, &session); err != nil {
			err = errors.Wrap(err, "error unmarshal event.data.raw into session")
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			var items []*entity.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)

			tr := otel.Tracer("rabbitmq")
			ctx, span := tr.Start(c.Request.Context(), fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
			defer span.End()

			_ = broker.PublishEvent(ctx, broker.PublishEventReq{
				Channel:  h.channel,
				Routing:  broker.FanOut,
				Queue:    "",
				Exchange: broker.EventOrderPaid,
				Body: &entity.Order{
					ID:          session.Metadata["orderID"],
					CustomerID:  session.Metadata["customerID"],
					Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
					PaymentLink: session.Metadata["paymentLink"],
					Items:       items,
				},
			})
		}
	}
	c.JSON(http.StatusOK, nil)
}
