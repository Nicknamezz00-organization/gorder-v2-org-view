package service

import (
	"context"

	"github.com/Nicknamezz00/gorder-v2/common/metrics"
	"github.com/Nicknamezz00/gorder-v2/stock/adapters"
	"github.com/Nicknamezz00/gorder-v2/stock/app"
	"github.com/Nicknamezz00/gorder-v2/stock/app/query"
	"github.com/Nicknamezz00/gorder-v2/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	stripeAPI := integration.NewStripeAPI()
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
