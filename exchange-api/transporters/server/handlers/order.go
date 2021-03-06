package handlers

import (
	"encoding/json"
	"exchange-api/services"
	"exchange-api/validators"
	"io/ioutil"
	"net/http"
)

type (
	OrderHandlerInterface interface {
		HandleOrder(w http.ResponseWriter, r *http.Request)
		CreateOrder(w http.ResponseWriter, r *http.Request)
		GetOrdersDashboardByCustomer(w http.ResponseWriter, r *http.Request)
	}

	OrderHandler struct {
		OrderService services.OrderServiceInterface
	}
)

func NewOrderHandler(s services.OrderServiceInterface) OrderHandlerInterface {
	return OrderHandler{
		OrderService: s,
	}
}

func (h OrderHandler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.CreateOrder(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))

	}
}

func (h OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	orderPayload := &services.CreateOrderSchema{}

	err = json.Unmarshal(body, orderPayload)

	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	createdOrder, err := h.OrderService.CreateOrder(ctx, orderPayload)

	if err != nil {
		if _, ok := err.(*validators.OrderValidationError); ok {

			response, _ := json.Marshal(map[string]interface{}{
				"message": err.(*validators.OrderValidationError).Message,
				"type":    "BadRequest",
				"details": err.(*validators.OrderValidationError).Details,
			})

			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else {

			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	response, err := json.Marshal(createdOrder)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h OrderHandler) GetOrdersDashboardByCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))

		return
	}

	ctx := r.Context()
	customerId := r.Header.Get("x-customer-id")

	if customerId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Customer id is required"))

		return
	}

	result, err := h.OrderService.GetOrdersDashboardByCustomer(ctx, customerId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(result)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
