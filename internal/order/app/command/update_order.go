package command

import (
	"context"

	"github.com/Nicknamezz00/gorder-v2/common/decorator"
	"github.com/Nicknamezz00/gorder-v2/common/logging"
	domain "github.com/Nicknamezz00/gorder-v2/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
	//stockGRPC
}

func NewUpdateOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) UpdateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	var err error
	defer logging.WhenCommandExecute(ctx, "UpdateOrderHandler", cmd, err)

	if cmd.UpdateFn == nil {
		logrus.Panicf("UpdateOrderHandler got nil order, cmd=%+v", cmd)
	}
	if err = c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
		return nil, err
	}
	return nil, nil
}
