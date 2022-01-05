package main

import (
	"exchange-api/repositories"
	"exchange-api/services"
	"exchange-api/transporters/server/handlers"
	"exchange-api/transporters/server/middlewares"
	"exchange-api/validators"
	"net/http"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func setupHandlers(h *http.ServeMux, sqs *sqs.SQS) http.Handler {
	validator := validators.NewValidator()
	orderRepository := repositories.NewOrderRepository(sqs)
	orderService := services.NewOrderService(orderRepository, validator)
	OrderHandler := handlers.NewOrderHandler(orderService)
	h.HandleFunc("/orders", OrderHandler.HandleOrder)

	h.HandleFunc("/__health_check__", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is up and kicking"))
	})

	return middlewares.HttpLogger(h)
}
