package handlers

import (
	"encoding/json"
	"exchange-api/services"
	"io/ioutil"
	"net/http"
)

type (
	ExchangeRateHandlerInterface interface {
		HandleExchangeRate(w http.ResponseWriter, r *http.Request)
		CreateExchangeRate(w http.ResponseWriter, r *http.Request)
	}

	ExchangeRateHandler struct {
		ExchangeRateService services.ExchangeRateServiceInterface
	}
)

func NewExchangeRateHandler(s services.ExchangeRateServiceInterface) ExchangeRateHandlerInterface {
	return ExchangeRateHandler{
		ExchangeRateService: s,
	}
}

func (h ExchangeRateHandler) HandleExchangeRate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.CreateExchangeRate(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))

	}
}

func (h ExchangeRateHandler) CreateExchangeRate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	payload := &services.CreateExchangeRateSchema{}

	err = json.Unmarshal(body, payload)

	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	createdExchangeRate, err := h.ExchangeRateService.CreateExchangeRate(ctx, payload)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	}

	response, err := json.Marshal(createdExchangeRate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
