package service

import (
	"context"

	"github.com/Nicknamezz00/gorder-v2/common/metrics"
	"github.com/Nicknamezz00/gorder-v2/order/adapters"
	"github.com/Nicknamezz00/gorder-v2/order/app"
	"github.com/Nicknamezz00/gorder-v2/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
