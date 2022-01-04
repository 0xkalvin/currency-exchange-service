package main

import (
	"exchange-api/repositories"
	"exchange-api/services"
	"exchange-api/transporters/server/handlers"
	"exchange-api/transporters/server/middlewares"
	"net/http"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func setupHandlers(h *http.ServeMux, sqs *sqs.SQS) http.Handler {
	orderRepository := repositories.NewOrderRepository(sqs)
	orderService := services.NewOrderService(orderRepository)
	OrderHandler := handlers.NewOrderHandler(orderService)
	h.HandleFunc("/orders", OrderHandler.HandleOrder)

	h.HandleFunc("/__health_check__", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is up and kicking"))
	})

	return middlewares.HttpLogger(h)
}
