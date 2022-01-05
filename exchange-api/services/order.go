package services

import (
	"context"
	"exchange-api/entities"
	"exchange-api/repositories"
	"exchange-api/utils/logger"
	"exchange-api/validators"

	"github.com/google/uuid"
)

type (
	OrderServiceInterface interface {
		CreateOrder(ctx context.Context, orderPayload *CreateOrderSchema) (*entities.Order, error)
	}

	OrderService struct {
		OrderRepository repositories.OrderRepositoryInterface
		Validator       validators.OrderValidatorInterface
	}

	CreateOrderSchema struct {
		Amount           string `json:"amount" validate:"required"`
		CustomerId       string `json:"customer_id" validate:"required"`
		SourceCurrencyId string `json:"source_currency_id" validate:"required"`
		TargetCurrencyId string `json:"target_currency_id" validate:"required"`
	}
)

var log = logger.NewLogger()

func NewOrderService(r repositories.OrderRepositoryInterface, v validators.OrderValidatorInterface) OrderServiceInterface {
	return OrderService{
		OrderRepository: r,
		Validator:       v,
	}
}

func (s OrderService) CreateOrder(ctx context.Context, orderPayload *CreateOrderSchema) (*entities.Order, error) {
	err := s.Validator.Validate(orderPayload)

	if err != nil {
		return nil, err
	}

	order, err := s.OrderRepository.CreateOrder(ctx, &entities.Order{
		Id:               uuid.New().String(),
		Amount:           orderPayload.Amount,
		CustomerId:       orderPayload.CustomerId,
		SourceCurrencyId: orderPayload.SourceCurrencyId,
		TargetCurrencyId: orderPayload.TargetCurrencyId,
	})

	if err != nil {
		return nil, err
	}

	log.Debug(map[string]interface{}{
		"message":    "Successfully created order",
		"order_id":   order.Id,
		"request_id": ctx.Value("requestId").(string),
	})

	return order, nil
}
