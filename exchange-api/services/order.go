package services

import (
	"exchange-api/entities"
	"exchange-api/repositories"
	"fmt"

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

func NewOrderService(r repositories.OrderRepositoryInterface) OrderServiceInterface {
	return OrderService{
		OrderRepository: r,
	}
}

func (s OrderService) CreateOrder(orderPayload *entities.Order) (*entities.Order, error) {
	orderPayload.Id = uuid.New().String()

	order, err := s.OrderRepository.CreateOrder(orderPayload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return order, nil
}
