package domain

import (
	"context"

	"github.com/Nicknamezz00/gorder-v2/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
