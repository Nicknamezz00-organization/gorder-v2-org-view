package ports

import (
	"context"

	"github.com/Nicknamezz00/gorder-v2/common/genproto/stockpb"
	"github.com/Nicknamezz00/gorder-v2/common/tracing"
	"github.com/Nicknamezz00/gorder-v2/stock/app"
	"github.com/Nicknamezz00/gorder-v2/stock/app/query"
	"github.com/Nicknamezz00/gorder-v2/stock/convertor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	_, span := tracing.Start(ctx, "GetItems")
	defer span.End()

	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIDs: request.ItemIDs})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &stockpb.GetItemsResponse{Items: convertor.NewItemConvertor().EntitiesToProtos(items)}, nil
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	_, span := tracing.Start(ctx, "CheckIfItemsInStock")
	defer span.End()

	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{
		Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &stockpb.CheckIfItemsInStockResponse{
		InStock: 1,
		Items:   convertor.NewItemConvertor().EntitiesToProtos(items),
	}, nil
}
