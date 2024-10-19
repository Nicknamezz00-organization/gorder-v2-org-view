package app

import "github.com/Nicknamezz00/gorder-v2/payment/app/command"

type Application struct {
	Commands Commands
}

type Commands struct {
	CreatePayment command.CreatePaymentHandler
}
