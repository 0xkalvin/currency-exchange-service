package handlers

import (
	"encoding/json"
	"exchange-api/services"
	"net/http"
)

type (
	BalanceHandlerInterface interface {
		GetMovementsDashboardByCustomer(w http.ResponseWriter, r *http.Request)
	}

	BalanceHandler struct {
		BalanceService services.BalanceServiceInterface
	}
)

func NewBalanceHandler(s services.BalanceServiceInterface) BalanceHandlerInterface {
	return BalanceHandler{
		BalanceService: s,
	}
}

func (h BalanceHandler) GetMovementsDashboardByCustomer(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.BalanceService.GetMovementsDashboardByCustomer(ctx, customerId)

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
