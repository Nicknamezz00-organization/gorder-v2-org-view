package command

import (
	"context"

	"github.com/Nicknamezz00/gorder-v2/common/decorator"
	"github.com/Nicknamezz00/gorder-v2/common/genproto/orderpb"
	"github.com/Nicknamezz00/gorder-v2/order/app/query"
	domain "github.com/Nicknamezz00/gorder-v2/order/domain/order"
	"github.com/sirupsen/logrus"
)

type CreateOrder struct {
	CustomerID string
	Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
		logger,
		metricClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// TODO: call stock grpc to get items.
	err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
	resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})
	logrus.Info("createOrderHandler||resp from stockGRPC.GetItems", resp)
	var stockResponse []*orderpb.Item
	for _, item := range cmd.Items {
		stockResponse = append(stockResponse, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      stockResponse,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{OrderID: o.ID}, nil
}
