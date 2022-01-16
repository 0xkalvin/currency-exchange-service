package services

import (
	"context"
	"exchange-api/repositories"
)

type (
	BalanceServiceInterface interface {
		GetMovementsDashboardByCustomer(ctx context.Context, customerId string) ([]*repositories.MovementTypeTotalItem, error)
		GetBalanceByCustomer(ctx context.Context, customerId string) ([]*repositories.BalanceResponseItem, error)
	}

	BalanceService struct {
		BalanceRepository repositories.BalanceRepositoryInterface
	}
)

func NewBalanceService(r repositories.BalanceRepositoryInterface) BalanceServiceInterface {
	return BalanceService{
		BalanceRepository: r,
	}
}

func (s BalanceService) GetMovementsDashboardByCustomer(ctx context.Context, customerId string) ([]*repositories.MovementTypeTotalItem, error) {
	return s.BalanceRepository.GetMovementsDashboardByCustomer(ctx, customerId)
}

func (s BalanceService) GetBalanceByCustomer(ctx context.Context, customerId string) ([]*repositories.BalanceResponseItem, error) {
	return s.BalanceRepository.GetBalanceByCustomer(ctx, customerId)
}
