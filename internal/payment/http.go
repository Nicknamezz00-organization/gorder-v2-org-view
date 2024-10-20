package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Nicknamezz00/gorder-v2/common/broker"
	"github.com/Nicknamezz00/gorder-v2/common/genproto/orderpb"
	"github.com/Nicknamezz00/gorder-v2/payment/domain"
	"github.com/gin-gonic/gin"
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

func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.handleWebhook)
}

func (h *PaymentHandler) handleWebhook(c *gin.Context) {
	logrus.Info("receive webhook from stripe")
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Infof("Error reading request body: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}

	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("ENDPOINT_STRIPE_SECRET"))

	if err != nil {
		logrus.Infof("Error verifying webhook signature: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			logrus.Infof("error unmarshal event.data.raw into session, err = %v", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Infof("payment for checkout session %v success!", session.ID)

			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			var items []*orderpb.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)

			marshalledOrder, err := json.Marshal(&domain.Order{
				ID:          session.Metadata["orderID"],
				CustomerID:  session.Metadata["customerID"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				PaymentLink: session.Metadata["paymentLink"],
				Items:       items,
			})
			if err != nil {
				logrus.Infof("error marshal domain.order, err = %v", err)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			tr := otel.Tracer("rabbitmq")
			mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
			defer span.End()

			headers := broker.InjectRabbitMQHeaders(mqCtx)
			_ = h.channel.PublishWithContext(mqCtx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         marshalledOrder,
				Headers:      headers,
			})
			logrus.Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
		}
	}
	c.JSON(http.StatusOK, nil)
}
