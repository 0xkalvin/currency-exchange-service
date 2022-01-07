package main

import (
	"exchange-api/repositories"
	"exchange-api/services"
	"exchange-api/transporters/server/handlers"
	"exchange-api/transporters/server/middlewares"
	"exchange-api/validators"
	"net/http"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func setupHandlers(h *http.ServeMux, sqs *sqs.SQS, dynamo *dynamodb.DynamoDB) http.Handler {
	validator := validators.NewValidator()
	orderRepository := repositories.NewOrderRepository(sqs, dynamo)
	orderService := services.NewOrderService(orderRepository, validator)
	OrderHandler := handlers.NewOrderHandler(orderService)

	exchangeRateRepository := repositories.NewExchangeRateRepository(dynamo)
	exchangeRateService := services.NewExchangeRateService(exchangeRateRepository)
	exchangeRateHandler := handlers.NewExchangeRateHandler(exchangeRateService)

	balanceRepository := repositories.NewBalanceRepository(dynamo)
	balanceService := services.NewBalanceService(balanceRepository)
	balanceHandler := handlers.NewBalanceHandler(balanceService)

	h.HandleFunc(
		"/orders",
		middlewares.AddRequestId(
			middlewares.HttpLogger(
				OrderHandler.HandleOrder),
		),
	)

	h.HandleFunc(
		"/orders/dashboard",
		middlewares.AddRequestId(
			middlewares.HttpLogger(
				OrderHandler.GetOrdersDashboardByCustomer),
		),
	)

	h.HandleFunc(
		"/exchange_rates",
		middlewares.AddRequestId(
			middlewares.HttpLogger(
				exchangeRateHandler.HandleExchangeRate),
		),
	)

	h.HandleFunc(
		"/movements/dashboard",
		middlewares.AddRequestId(
			middlewares.HttpLogger(
				balanceHandler.GetMovementsDashboardByCustomer),
		),
	)

	h.HandleFunc("/__health_check__", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is up and kicking"))
	})

	return h
}
