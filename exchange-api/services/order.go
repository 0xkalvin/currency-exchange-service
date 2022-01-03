package services

import (
	"exchange-api/entities"
	"exchange-api/repositories"
	"exchange-api/utils/logger"

	"github.com/google/uuid"
)

type (
	OrderServiceInterface interface {
		CreateOrder(orderPayload *entities.Order) (*entities.Order, error)
	}

	OrderService struct {
		OrderRepository repositories.OrderRepositoryInterface
	}
)

var log = logger.NewLogger()

func NewOrderService(r repositories.OrderRepositoryInterface) OrderServiceInterface {
	return OrderService{
		OrderRepository: r,
	}
}

func (s OrderService) CreateOrder(orderPayload *entities.Order) (*entities.Order, error) {
	orderPayload.Id = uuid.New().String()

	order, err := s.OrderRepository.CreateOrder(orderPayload)

	if err != nil {
		return nil, err
	}

	log.Debug(map[string]interface{}{
		"message":  "Successfully created order",
		"order_id": order.Id,
	})

	return order, nil
}
