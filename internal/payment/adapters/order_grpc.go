package adapters

import (
	"context"

	"github.com/Nicknamezz00/gorder-v2/common/genproto/orderpb"
	"github.com/Nicknamezz00/gorder-v2/common/tracing"
	"github.com/sirupsen/logrus"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
	ctx, span := tracing.Start(ctx, "order_grpc.update_order")
	defer span.End()

	_, err := o.client.UpdateOrder(ctx, order)
	logrus.Infof("payment_adapter||update_order,err=%v", err)
	return err
}
