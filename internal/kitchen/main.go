package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nicknamezz00/gorder-v2/common/broker"
	grpcClient "github.com/Nicknamezz00/gorder-v2/common/client"
	_ "github.com/Nicknamezz00/gorder-v2/common/config"
	"github.com/Nicknamezz00/gorder-v2/common/logging"
	"github.com/Nicknamezz00/gorder-v2/common/tracing"
	"github.com/Nicknamezz00/gorder-v2/kitchen/adapters"
	"github.com/Nicknamezz00/gorder-v2/kitchen/infrastructure/consumer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("kitchen.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	orderClient, closeFunc, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer closeFunc()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	orderGRPC := adapters.NewOrderGRPC(orderClient)
	go consumer.NewConsumer(orderGRPC).Listen(ch)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		logrus.Infof("receive signal, exiting...")
		os.Exit(0)
	}()
	logrus.Println("to exit, press ctrl+c")
	select {}
}
