package handlers

import (
	"encoding/json"
	"exchange-api/entities"
	"exchange-api/services"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	OrderHandlerInterface interface {
		HandleOrder(w http.ResponseWriter, r *http.Request)
		CreateOrder(w http.ResponseWriter, r *http.Request)
	}

	OrderHandler struct {
		OrderService services.OrderServiceInterface
	}

	CreateOrderSchema struct {
		Amount           string `json:"amount"`
		CustomerId       string `json:"customer_id"`
		SourceCurrencyId string `json:"source_currency_id"`
		TargetCurrencyId string `json:"target_currency_id"`
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
		fmt.Println(r.Method)
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

	orderPayload := CreateOrderSchema{}

	err = json.Unmarshal(body, &orderPayload)

	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)

		return
	}

	createdOrder, err := h.OrderService.CreateOrder(&entities.Order{
		Amount:           orderPayload.Amount,
		CustomerId:       orderPayload.CustomerId,
		SourceCurrencyId: orderPayload.SourceCurrencyId,
		TargetCurrencyId: orderPayload.TargetCurrencyId,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

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
