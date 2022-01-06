package services

import (
	"context"
	"exchange-api/entities"
	"exchange-api/repositories"
	"exchange-api/utils/logger"
	"time"

	"github.com/google/uuid"
)

type (
	ExchangeRateServiceInterface interface {
		CreateExchangeRate(ctx context.Context, payload *CreateExchangeRateSchema) (*entities.ExchangeRate, error)
	}

	ExchangeRateService struct {
		ExchangeRateRepository repositories.ExchangeRateRepositoryInterface
	}

	CreateExchangeRateSchema struct {
		CurrencyId string `json:"currency_id" validate:"required"`
		Rate       string `json:"rate" validate:"required"`
	}
)

const baseCurrencyId = "USD"

var logExchangeRate = logger.NewLogger()

func NewExchangeRateService(r repositories.ExchangeRateRepositoryInterface) ExchangeRateServiceInterface {
	return ExchangeRateService{
		ExchangeRateRepository: r,
	}
}

func (s ExchangeRateService) CreateExchangeRate(ctx context.Context, payload *CreateExchangeRateSchema) (*entities.ExchangeRate, error) {

	exchangeRate, err := s.ExchangeRateRepository.CreateExchangeRate(ctx, &entities.ExchangeRate{
		Id:             uuid.New().String(),
		CurrencyId:     payload.CurrencyId,
		BaseCurrencyId: baseCurrencyId,
		Rate:           payload.Rate,
		Timestamp:      time.Now().UTC().String(),
	})

	if err != nil {
		return nil, err
	}

	logExchangeRate.Debug(map[string]interface{}{
		"message":          "Successfully created exchange rate",
		"exchange_rate_id": exchangeRate.Id,
		"request_id":       ctx.Value("requestId").(string),
	})

	return exchangeRate, nil
}
