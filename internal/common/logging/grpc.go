package logging

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GRPCUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	fields := logrus.Fields{
		Args: req,
	}
	defer func() {
		fields[Response] = resp
		if err != nil {
			fields[Error] = err.Error()
			logf(ctx, logrus.ErrorLevel, fields, "%s", "_grpc_request_out")
		}
	}()
	md, exist := metadata.FromIncomingContext(ctx)
	if exist {
		fields["grpc_metadata"] = md
	}

	logf(ctx, logrus.InfoLevel, fields, "%s", "_grpc_request_in")
	return handler(ctx, req)
}
