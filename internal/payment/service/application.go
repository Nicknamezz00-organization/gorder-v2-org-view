package service

import (
	"context"

	grpcClient "github.com/Nicknamezz00/gorder-v2/common/client"
	"github.com/Nicknamezz00/gorder-v2/common/metrics"
	"github.com/Nicknamezz00/gorder-v2/payment/adapters"
	"github.com/Nicknamezz00/gorder-v2/payment/app"
	"github.com/Nicknamezz00/gorder-v2/payment/app/command"
	"github.com/Nicknamezz00/gorder-v2/payment/domain"
	"github.com/Nicknamezz00/gorder-v2/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	memoryProcessor := processor.NewInmemProcessor()
	return newApplication(ctx, orderGRPC, memoryProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(ctx context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}
